package main

import (
	"flag"
	"fmt"
	"os"
)

var version string

type Flags struct {
	Config string
	Dryrun bool
	Debug  bool
}

func parseFlags() (flags *Flags) {
	flags = &Flags{}
	flag.StringVar(&flags.Config, "config", "", "Config file path")
	flag.BoolVar(&flags.Dryrun, "dryrun", false, "Dry-run mode")
	flag.BoolVar(&flags.Debug, "debug", false, "Debug mode")
	flagVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if flags.Config == "" {
		printErrorAndExit("'-config' is required")
	}

	if *flagVersion {
		printVersionAndEixt()
	}

	return
}

func printVersionAndEixt() {
	fmt.Fprintln(os.Stderr, version)
	os.Exit(0)
}

func printErrorAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
