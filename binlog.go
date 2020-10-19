package binrpt

import (
	"context"
	"fmt"

	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

type Binlog struct {
	*SourceConfig
}

func NewBinlog(config *SourceConfig) *Binlog {
	return &Binlog{config}
}

func (binlog *Binlog) Receive(evout chan Event, ctx context.Context) error {
	file, pos, err := binlog.sourceStatus()

	if err != nil {
		return fmt.Errorf("Failed to get binlog status: %w", err)
	}

	return binlog.startSync(file, pos, evout, ctx)
}

func (binlog *Binlog) sourceStatus() (file string, pos uint32, err error) {
	conn, err := binlog.Connect()

	if err != nil {
		return
	}

	defer conn.Close()
	r, err := conn.Execute("show master status")

	if err != nil {
		return
	}

	file, err = r.GetString(0, 0)

	if err != nil {
		return
	}

	pos64, err := r.GetUint(0, 1)
	pos = uint32(pos64)

	return
}

func (binlog *Binlog) startSync(file string, pos uint32, evout chan Event, ctx context.Context) error {
	syncer, err := binlog.NewBinlogSyncer(binlog.ReplicateServerId)

	if err != nil {
		return err
	}

	defer syncer.Close()
	streamer, err := syncer.StartSync(mysql.Position{Name: file, Pos: pos})

	if err != nil {
		return err
	}

	for {
		ev, err := streamer.GetEvent(ctx)

		if err != nil {
			return err
		}

		switch ev.Header.EventType {
		case replication.WRITE_ROWS_EVENTv2:
			binlog.handleRowsEvent(ev, file, evout)
		case replication.UPDATE_ROWS_EVENTv2:
			binlog.handleRowsEvent(ev, file, evout)
		case replication.DELETE_ROWS_EVENTv2:
			binlog.handleRowsEvent(ev, file, evout)
		case replication.QUERY_EVENT:
			binlog.handleQueryEvent(ev, file, evout)
		case replication.ROTATE_EVENT:
			rotateEvent := ev.Event.(*replication.RotateEvent)
			file = string(rotateEvent.NextLogName)
		}
	}
}

func (binlog *Binlog) handleRowsEvent(ev *replication.BinlogEvent, file string, evout chan Event) {
	event := Event{
		File:      file,
		Header:    ev.Header,
		RowsEvent: ev.Event.(*replication.RowsEvent),
	}

	evout <- event
}

func (binlog *Binlog) handleQueryEvent(ev *replication.BinlogEvent, file string, evout chan Event) {
	event := Event{
		File:       file,
		Header:     ev.Header,
		QueryEvent: ev.Event.(*replication.QueryEvent),
	}

	evout <- event
}
