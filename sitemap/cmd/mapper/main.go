package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Colocasian/gophercises/sitemap"
)

func main() {
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	log.SetFlags(0)

	flag.Parse()
	if args.url == "" {
		log.Fatalln("did not specify -url parameter")
	}

	sm, err := sitemap.CrawlMap(args.url, args.depth)
	if err != nil {
		log.Fatalf("could not create sitemap: %v", err)
	}

	txt, err := xml.MarshalIndent(sm, "", "   ")
	if err != nil {
		log.Fatalf("could not generate XML: %v", err)
	}

	if args.out == "" {
		fmt.Println(string(txt))
	} else {
		err := os.WriteFile(args.out, txt, 0644)
		if err != nil {
			log.Fatalf("could not write to file %q: %v", args.out, err)
		}
	}
}
