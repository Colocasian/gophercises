package main

import (
	"flag"
	"time"
)

type Cli struct {
	csv   string
	limit time.Duration
	shuf  bool
}

var cli Cli

func init() {
	flag.StringVar(&cli.csv, "csv", "problems.csv", "a csv file in the format of 'question, answer'")
	flag.DurationVar(&cli.limit, "limit", 30*time.Second, "the time limit for the quiz")
	flag.BoolVar(&cli.shuf, "shuf", false, "shuffle the question order")
}
