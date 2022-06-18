package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Set up logging for main application
	log.SetPrefix(fmt.Sprintf("%v: ", os.Args[0]))
	log.SetFlags(0)

	// Parse all CLI flags
	flag.Parse()

	file, err := os.Open(cli.csv)
	if err != nil {
		log.Fatalf("could not open file %q: %v\n", cli.csv, err)
	}

	quizzer := readCSV(file)
	score := quizzer.Conduct(cli.limit, cli.shuf)
	fmt.Printf("You scored %v out of %v.\n", score, quizzer.Length())
}
