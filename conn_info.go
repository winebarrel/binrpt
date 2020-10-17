package binrpt

import (
	"fmt"
	"os"
	"strconv"

	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/replication"
)

type ConnInfo struct {
	Host     string
	Port     uint16
	Username string
	Password string
	Charset  string
}

func (connInfo *ConnInfo) Connect() (*client.Conn, error) {
	hostPort := fmt.Sprintf("%s:%d", connInfo.Host, connInfo.Port)
	conn, err := client.Connect(hostPort, connInfo.Username, connInfo.Password, "")

	if err != nil {
		return nil, err
	}

	err = conn.SetCharset(connInfo.Charset)

	return conn, err
}

func (connInfo *ConnInfo) Ping() error {
	conn, err := connInfo.Connect()

	if err != nil {
		return err
	}

	defer conn.Close()
	return conn.Ping()
}

func (connInfo *ConnInfo) NewBinlogSyncer(serverId uint32) (*replication.BinlogSyncer, error) {
	maxReconnStr := os.Getenv("BINLOG_MAX_RECONNECT_ATTEMPTS")
	maxReconn := 0
	var err error

	if maxReconnStr != "" {
		maxReconn, err = strconv.Atoi(maxReconnStr)

		if err != nil {
			return nil, fmt.Errorf("BINLOG_MAX_RECONNECT_ATTEMPTS env parse failed: %w", err)
		}
	}

	cfg := replication.BinlogSyncerConfig{
		ServerID:             serverId,
		Flavor:               "mysql",
		Host:                 connInfo.Host,
		Port:                 connInfo.Port,
		User:                 connInfo.Username,
		Password:             connInfo.Password,
		Charset:              connInfo.Charset,
		MaxReconnectAttempts: maxReconn,
	}

	return replication.NewBinlogSyncer(cfg), nil
}
