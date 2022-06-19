package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
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

	rt, at, err := getTemplates()
	if err != nil {
		log.Fatalf("could not read HTML templates: %v", err)
	}

	handler := s.HTTPHandler(rt, at)
	log.Printf("Listening to :8080")
	http.ListenAndServe(":8080", handler)
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

func getTemplates() (*template.Template, *template.Template, error) {
	rt, err := template.ParseFiles("views/index.html")
	if err != nil {
		return nil, nil, err
	}
	at, err := template.ParseFiles("views/arc.html")
	if err != nil {
		return nil, nil, err
	}

	return rt, at, nil
}
