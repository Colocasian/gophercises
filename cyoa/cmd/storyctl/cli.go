package main

import "flag"

type Args struct {
	book string
}

var args Args

func init() {
	flag.StringVar(&args.book, "book", "", "file with the story in JSON format")
}
