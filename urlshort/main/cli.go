package main

import "flag"

type Cli struct {
	conf   string
	format string
}

var cli Cli

func init() {
	flag.StringVar(&cli.conf, "conf", "configs/config.db", "path of config file")
	flag.StringVar(&cli.format, "format", "bolt", "config file format")
}
