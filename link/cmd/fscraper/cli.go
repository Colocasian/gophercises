package main

import "flag"

type Args struct {
	file string
}

var args Args

func init() {
	flag.StringVar(&args.file, "file", "", "HTML file to parse")
}
