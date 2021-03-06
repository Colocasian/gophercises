package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Colocasian/gophercises/link"
	"gopkg.in/yaml.v3"
)

func main() {
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	log.SetFlags(0)

	flag.Parse()

	var (
		f   *os.File
		err error
	)
	if args.file == "" {
		f = os.Stdin
	} else {
		f, err = os.Open(args.file)
		if err != nil {
			log.Fatalln(err)
		}
	}

	links, err := link.ParseLinks(f)
	if err != nil {
		log.Fatalf("could not parse links: %v", err)
	}

	bytes, err := yaml.Marshal(links)
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(bytes))
}
