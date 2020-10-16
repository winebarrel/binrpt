package binrpt

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

const (
	DefaultPort = 3306
)

type MasterConfig struct {
	ConnInfo
	ServerId uint32 `toml:"server_id"`
}

type ReplicaConfig struct {
	ConnInfo
	ReplicateDoDB string `toml:"replicate_do_db"`
}

type Filter struct {
	IgnoreTable string `toml:"ignore_table"`
}

type Config struct {
	Master  MasterConfig
	Replica ReplicaConfig
	Filters []Filter
}

func LoadConfig(path string) (config *Config, err error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	config = &Config{}
	err = toml.Unmarshal(content, config)

	if err != nil {
		return nil, err
	}

	if config.Master.Host == "" {
		return nil, fmt.Errorf("config error: master.host is required")
	}

	if config.Master.Username == "" {
		return nil, fmt.Errorf("config error: master.username is required")
	}

	if config.Master.ServerId <= 0 {
		return nil, fmt.Errorf("config error: master.server_id mult be '>= 1'")
	}

	if config.Master.Port == 0 {
		config.Master.Port = 3306
	}

	if config.Replica.Host == "" {
		return nil, fmt.Errorf("config error: replica.host is required")
	}

	if config.Replica.Username == "" {
		return nil, fmt.Errorf("config error: replica.username is required")
	}

	if config.Replica.ReplicateDoDB == "" {
		return nil, fmt.Errorf("config error: replica.replicate_do_db is required")
	}

	if config.Replica.Port == 0 {
		config.Replica.Port = 3306
	}

	return config, nil
}
