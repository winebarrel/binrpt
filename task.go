package binrpt

import (
	"context"
	"fmt"
	"strings"

	"github.com/siddontang/go-log/log"

	"golang.org/x/sync/errgroup"
)

const (
	// https://github.com/siddontang/go-mysql/blob/v1.1.0/replication/row_event.go#L19
	ErrMissingTableMapEventMessage = "invalid table id, no corresponding table map event"
)

type Task struct {
	Config *Config
}

func NewTask(config *Config) *Task {
	return &Task{Config: config}
}

func (task *Task) Run(dryrun bool) error {
	binlog := NewBinlog(&task.Config.Source)
	err := binlog.Ping()

	if err != nil {
		return fmt.Errorf("Failed to create Binlog instance: %w", err)
	}

	replica, err := NewReplica(&task.Config.Replica, dryrun)

	if err != nil {
		return fmt.Errorf("Failed to create Replica instance: %w", err)
	}

	err = replica.Ping()

	if err != nil {
		return fmt.Errorf("Replica did not respond to pin: %w", err)
	}

	var file string
	var pos uint32

	if replica.SaveStatus {
		file, pos, err = replica.LoadBinlogFilePos()

		if err != nil {
			return fmt.Errorf("Failed to load binlog position: %w", err)
		}
	}

	evch := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())
	eg, etCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer func() {
			cancel()

			// Skip events
			for ev := range evch {
				_ = ev
			}
		}()

		err := replica.Repeat(evch, etCtx, binlog)

		if err != nil {
			log.Errorf("Faild to repeat binlog: %s", err)
		}

		return err
	})

	eg.Go(func() error {
		defer close(evch)
		err := binlog.Receive(evch, etCtx, file, pos)

		if strings.Contains(err.Error(), ErrMissingTableMapEventMessage) {
			file, pos, loadErr := replica.LoadBinlogMapEventFilePos()

			if loadErr != nil {
				return loadErr
			}
			fmt.Println(file, pos)

			if file != "" {
				err = binlog.Receive(evch, etCtx, file, pos)
			}
		}

		return err
	})

	return eg.Wait()
}
