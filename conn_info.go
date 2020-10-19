package binrpt

import (
	"fmt"

	"github.com/siddontang/go-mysql/client"
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
