package binrpt

import "github.com/siddontang/go-mysql/replication"

type Event struct {
	File          string
	Header        *replication.EventHeader
	RowsEvent     *replication.RowsEvent
	QueryEvent    *replication.QueryEvent
	TableMapEvent *replication.TableMapEvent
}
