package main

import "flag"

type Cli struct {
	conf   string
	format string
}

var cli Cli

func init() {
	flag.StringVar(&cli.conf, "conf", "config.yaml", "path of config file")
	flag.StringVar(&cli.format, "format", "yaml", "config file format")
}
