package binrpt

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

const (
	DefaultPort = 3306
)

type SourceConfig struct {
	ConnInfo
	ReplicateServerId uint32 `toml:"replicate_server_id"`
}

type ReplicaConfig struct {
	ConnInfo
	ReplicateDoDB         string   `toml:"replicate_do_db"`
	ReplicateIgnoreTables []string `toml:"replicate_ignore_tables"`
}

type Config struct {
	Source  SourceConfig
	Replica ReplicaConfig
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

	if config.Source.Host == "" {
		return nil, fmt.Errorf("config error: source.host is required")
	}

	if config.Source.Username == "" {
		return nil, fmt.Errorf("config error: source.username is required")
	}

	if config.Source.Charset == "" {
		return nil, fmt.Errorf("config error: source.charset is required")
	}

	if config.Source.ReplicateServerId <= 0 {
		return nil, fmt.Errorf("config error: source.replicate_server_id mult be '>= 1'")
	}

	if config.Source.Port == 0 {
		config.Source.Port = 3306
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

	if config.Replica.Charset == "" {
		return nil, fmt.Errorf("config error: replica.charset is required")
	}

	if config.Replica.Port == 0 {
		config.Replica.Port = 3306
	}

	return config, nil
}
