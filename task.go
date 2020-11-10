package binrpt

import (
	"context"
	"fmt"

	"github.com/siddontang/go-log/log"

	"golang.org/x/sync/errgroup"
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

		err := replica.Repeat(evch, etCtx)

		if err != nil {
			log.Errorf("Faild to repeat binlog: %s", err)
		}

		return err
	})

	eg.Go(func() error {
		defer close(evch)
		return binlog.Receive(evch, etCtx)
	})

	return eg.Wait()
}
