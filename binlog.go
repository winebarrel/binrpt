package binrpt

import (
	"context"
	"fmt"

	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

type TrashScanner struct{}

func (TrashScanner) Scan(interface{}) error {
	return nil
}

type Binlog struct {
	*SourceConfig
}

func NewBinlog(config *SourceConfig) *Binlog {
	return &Binlog{config}
}

func (binlog *Binlog) Receive(evout chan Event, ctx context.Context, file string, pos uint32) error {
	var err error

	if file == "" {
		if binlog.BinlogBufferNum > 0 {
			file, pos, err = binlog.sourcePrevBinlogLast(int(binlog.BinlogBufferNum))
		} else {
			file, pos, err = binlog.sourceStatus()
		}

		if err != nil {
			return fmt.Errorf("Failed to get binlog status: %w", err)
		}
	}

	return binlog.startSync(file, pos, evout, ctx)
}

func (binlog *Binlog) sourceStatus() (file string, pos uint32, err error) {
	conn, err := binlog.Connect()

	if err != nil {
		return
	}

	defer conn.Close()
	rows, err := conn.Query("show master status")

	if err != nil {
		return
	}

	defer rows.Close()
	columns, err := rows.Columns()

	if err != nil {
		return
	}

	colLen := len(columns)
	dest := make([]interface{}, colLen)
	dest[0] = &file
	dest[1] = &pos

	for i := 2; i < colLen; i++ {
		dest[i] = TrashScanner{}
	}

	rows.Next()
	err = rows.Scan(dest...)

	if err != nil {
		return
	}

	return
}

type binlogFilePos struct {
	File string
	Pos  uint32
}

func (binlog *Binlog) sourcePrevBinlogLast(bufNum int) (file string, pos uint32, err error) {
	conn, err := binlog.Connect()

	if err != nil {
		return
	}

	defer conn.Close()
	rows, err := conn.Query("show master logs")

	if err != nil {
		return
	}

	defer rows.Close()
	columns, err := rows.Columns()

	if err != nil {
		return
	}

	filePosList := []binlogFilePos{}

	var curFile string
	var curPos uint32

	colLen := len(columns)
	dest := make([]interface{}, colLen)
	dest[0] = &curFile
	dest[1] = &curPos

	for i := 2; i < colLen; i++ {
		dest[i] = TrashScanner{}
	}

	for rows.Next() {
		err = rows.Scan(dest...)

		if err != nil {
			return
		}

		filePos := binlogFilePos{File: curFile, Pos: curPos}
		filePosList = append(filePosList, filePos)
	}

	filePosListLen := len(filePosList)
	fmt.Println(filePosList)

	if filePosListLen <= bufNum {
		err = fmt.Errorf("Failed to get previous binlog position")
	}

	startFilePos := filePosList[filePosListLen-(1+bufNum)]
	file = startFilePos.File
	pos = startFilePos.Pos

	return
}

func (binlog *Binlog) startSync(file string, pos uint32, evout chan Event, ctx context.Context) error {
	syncer := binlog.NewBinlogSyncer(binlog.ReplicateServerId)
	defer syncer.Close()
	streamer, err := syncer.StartSync(mysql.Position{Name: file, Pos: pos})

	if err != nil {
		return fmt.Errorf("Failed to start syncing binlog: %w", err)
	}

	for {
		ev, err := streamer.GetEvent(ctx)

		if err != nil {
			return fmt.Errorf("Failed to get binlog event: %w", err)
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
