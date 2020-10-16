package main

import (
	"github.com/siddontang/go-log/log"
	"github.com/winebarrel/binrpt"
)

func main() {
	flags := parseFlags()
	config, err := binrpt.LoadConfig(flags.Config)

	if err != nil {
		log.Fatal(err)
	}

	task := binrpt.NewTask(config)
	err = task.Run()

	if err != nil {
		log.Fatal(err)
	}
}
