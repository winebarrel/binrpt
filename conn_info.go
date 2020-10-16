package binrpt

import (
	"fmt"

	"github.com/siddontang/go-mysql/client"
)

type ConnInfo struct {
	Host     string
	Port     uint16
	Username string
	Password string
}

func (connInfo *ConnInfo) Connect() (*client.Conn, error) {
	hostPort := fmt.Sprintf("%s:%d", connInfo.Host, connInfo.Port)
	return client.Connect(hostPort, connInfo.Username, connInfo.Password, "")
}

func (connInfo *ConnInfo) Ping() error {
	conn, err := connInfo.Connect()

	if err != nil {
		return err
	}

	defer conn.Close()
	return conn.Ping()
}
