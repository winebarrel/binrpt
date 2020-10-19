package binrpt

import (
	"context"
	"regexp"
	"time"

	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/replication"
)

const (
	LogPosInterval    = 1 * time.Minute
	ReconnectInterval = 1 * time.Second
)

var ignoreQuery = regexp.MustCompile(`(?i)^(begin$|grant|revoke|create\s+user|drop\s+user)`)

type Replica struct {
	*ReplicaConfig
	IgnoreTableFilters []*regexp.Regexp
	Dryrun             bool
}

func NewReplica(config *ReplicaConfig, dryrun bool) (*Replica, error) {
	ignoreTableFilters := make([]*regexp.Regexp, len(config.ReplicateIgnoreTables))
	var err error

	for i, v := range config.ReplicateIgnoreTables {
		ignoreTableFilters[i], err = regexp.Compile(v)

		if err != nil {
			return nil, err
		}
	}

	return &Replica{
		ReplicaConfig:      config,
		IgnoreTableFilters: ignoreTableFilters,
		Dryrun:             dryrun,
	}, nil
}

func (replica *Replica) Repeat(evin chan Event, ctx context.Context) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	conn, err := replica.Connect()

	if err != nil {
		return err
	}

	defer conn.Close()
	tableInfo := NewTableInfo(conn)

	ticker := time.NewTicker(LogPosInterval)
	defer ticker.Stop()

	for ev := range evin {
		var err error
		conn, err = replica.pingAndReconnect(conn)

		if err != nil {
			return err
		}

		if ev.RowsEvent != nil {
			err = replica.handleRowsEvent(conn, ev.Header, ev.RowsEvent, tableInfo)
		} else if ev.QueryEvent != nil {
			err = replica.handleQueryEvent(conn, ev.Header, ev.QueryEvent)
		}

		if err != nil {
			return err
		}

		select {
		case <-ticker.C:
			log.Infof("log_file=%s log_pos=%d", ev.File, ev.Header.LogPos)
		default:
		}
	}

	return nil
}

func (replica *Replica) pingAndReconnect(conn *client.Conn) (*client.Conn, error) {
	err := conn.Ping()

	if err == nil {
		return conn, nil
	}

	log.Warnf("reconnect attempt: %s", err)
	reconnCount := 0

	for {
		if replica.MaxReconnectAttempts > 0 {
			if reconnCount >= replica.MaxReconnectAttempts {
				break
			}

			reconnCount++
		}

		conn, err = replica.Connect()

		if err == nil {
			return conn, nil
		}

		log.Warnf("reconnect attempt: %s", err)
		time.Sleep(ReconnectInterval)
	}

	return nil, err
}

func (replica *Replica) handleRowsEvent(conn *client.Conn, header *replication.EventHeader, ev *replication.RowsEvent, tableInfo *TableInfo) error {
	schema := string(ev.Table.Schema)

	if schema != replica.ReplicateDoDB {
		return nil
	}

	tableName := string(ev.Table.Table)

	for _, re := range replica.IgnoreTableFilters {
		if re.MatchString(tableName) {
			return nil
		}
	}

	table, err := tableInfo.Get(schema, tableName)

	if err != nil {
		return err
	}

	if table == nil {
		log.Warnf("Table not found: %s.%s", schema, tableName)
		return nil
	}

	if table.ColumnCount < ev.ColumnCount {
		log.Warnf("Table column count is less than ROWS_EVENT column count (%d < %d): %s.%s", table.ColumnCount, ev.ColumnCount, schema, tableName)
		return nil
	}

	sqlBld := NewSQLBuilder(table, header, ev)
	sqls := sqlBld.SQLs()

	for _, v := range sqls {
		if !replica.Dryrun {
			if len(v.Params) > 0 {
				_, err = conn.Execute(v.Statement, v.Params...)
			} else {
				_, err = conn.Execute(v.Statement)
			}

			if err != nil {
				log.Warnf("%s: %v", err, v)
			}
		} else {
			log.Infof("dry-run: %s %s", v.Statement, v.Params)
		}
	}

	return nil
}

func (replica *Replica) handleQueryEvent(conn *client.Conn, header *replication.EventHeader, ev *replication.QueryEvent) error {
	schema := string(ev.Schema)

	if schema != replica.ReplicateDoDB {
		return nil
	}

	query := string(ev.Query)

	if ignoreQuery.MatchString(query) {
		return nil
	}

	useStmt := "USE " + schema

	if !replica.Dryrun {
		_, err := conn.Execute(useStmt)

		if err != nil {
			log.Warnf("%s: %s", err, useStmt)
			return nil
		}

		_, err = conn.Execute(query)

		if err != nil {
			log.Warnf("%s: %s", err, query)
		}
	} else {
		log.Infof("dry-run: %s", useStmt)
		log.Infof("dry-run: %s", query)
	}

	return nil
}
