package binrpt

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Task struct {
	Config *Config
}

func NewTask(config *Config) *Task {
	return &Task{Config: config}
}

func (task *Task) Run() error {
	binlog := NewBinlog(&task.Config.Master)
	err := binlog.Ping()

	if err != nil {
		return err
	}

	replica, err := NewReplica(&task.Config.Replica)

	if err != nil {
		return err
	}

	err = replica.Ping()

	if err != nil {
		return err
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
