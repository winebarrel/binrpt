package main

import (
	"fmt"

	"github.com/siddontang/go-log/log"
	"github.com/winebarrel/binrpt"
)

func main() {
	flags := parseFlags()
	config, err := binrpt.LoadConfig(flags.Config)

	if err != nil {
		log.Fatal(fmt.Errorf("Failed to load config: %w", err))
	}

	task := binrpt.NewTask(config)
	err = task.Run(flags.Dryrun)

	if err != nil {
		log.Fatal(err)
	}
}
