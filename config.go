package binrpt

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/pelletier/go-toml"
)

const (
	DefaultPort       = 3306
	BinlogStatusDB    = "binrpt"
	BinlogStatusTable = "replica_status"
)

type SourceConfig struct {
	ConnInfo
	ReplicateServerId uint32 `toml:"replicate_server_id"`
	BinlogBufferNum   uint32 `toml:"binlog_buffer_num"`
}

type ReplicaConfig struct {
	ConnInfo
	ReplicateDoDB         string   `toml:"replicate_do_db"`
	ReplicateIgnoreTables []string `toml:"replicate_ignore_tables"`
	SaveStatus            bool     `toml:"save_status"`
	TableInfoFromSrc      bool     `toml:"table_info_from_src"`
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

	maxReconnStr := os.Getenv("REPLICATE_MAX_RECONNECT_ATTEMPTS")

	if maxReconnStr != "" {
		maxReconn, err := strconv.Atoi(maxReconnStr)

		if err != nil {
			return nil, fmt.Errorf("REPLICATE_MAX_RECONNECT_ATTEMPTS env parse failed: %w", err)
		}

		config.Source.MaxReconnectAttempts = maxReconn
		config.Replica.MaxReconnectAttempts = maxReconn
	}

	return config, nil
}
