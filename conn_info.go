package binrpt

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/siddontang/go-mysql/replication"
)

type ConnInfo struct {
	Host                 string
	Port                 uint16
	Username             string
	Password             string
	Charset              string
	MaxReconnectAttempts int
}

func (connInfo *ConnInfo) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s",
		connInfo.Username, connInfo.Password, connInfo.Host, connInfo.Port, connInfo.Charset)
	conn, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	conn.SetConnMaxLifetime(1 * time.Hour)
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	return conn, nil
}

func (connInfo *ConnInfo) Ping() error {
	conn, err := connInfo.Connect()

	if err != nil {
		return err
	}

	defer conn.Close()
	return conn.Ping()
}

func (connInfo *ConnInfo) NewBinlogSyncer(serverId uint32) *replication.BinlogSyncer {
	cfg := replication.BinlogSyncerConfig{
		ServerID:             serverId,
		Flavor:               "mysql",
		Host:                 connInfo.Host,
		Port:                 connInfo.Port,
		User:                 connInfo.Username,
		Password:             connInfo.Password,
		Charset:              connInfo.Charset,
		MaxReconnectAttempts: connInfo.MaxReconnectAttempts,
	}

	return replication.NewBinlogSyncer(cfg)
}
