package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Colocasian/gophercises/cyoa/pkg/story"
)

func main() {
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	log.SetFlags(0)
	flag.Parse()

	s, err := parseStory(args.book)
	if err != nil {
		log.Fatalf("could not get story: %v", err)
	}

	if err = s.StartCLI("intro"); err != nil {
		log.Fatalf("could not start story: %v", err)
	}
}

func parseStory(book string) (story.Story, error) {
	if book == "" {
		return nil, fmt.Errorf("did not specify book")
	}

	b, err := os.ReadFile(book)
	if err != nil {
		return nil, err
	}

	s, err := story.ParseStory(b)
	if err != nil {
		return nil, err
	}

	return s, nil
}
