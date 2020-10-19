package binrpt

import (
	"context"
	"fmt"

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
		return err
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
	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return replica.Repeat(evch, ctx)
	})

	eg.Go(func() error {
		defer close(evch)
		return binlog.Receive(evch, ctx)
	})

	return eg.Wait()
}
