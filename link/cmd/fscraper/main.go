package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Colocasian/gophercises/link/pkg/link"
)

func main() {
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	log.SetFlags(0)

	flag.Parse()
	if args.file == "" {
		log.Fatalln("did not specify HTML file to parse")
	}

	f, err := os.Open(args.file)
	if err != nil {
		log.Fatalln(err)
	}

	links, err := link.ParseLinks(f)
	if err != nil {
		log.Fatalf("could not parse links: %v", err)
	}

	fmt.Println(links)
}
