package main

import "flag"

var args struct {
	url   string
	depth int
	out string
}

func init() {
	flag.StringVar(&args.url, "url", "", "root URL to start sitemap from")
	flag.IntVar(&args.depth, "depth", 4, "depth of crawling")
	flag.StringVar(&args.out, "out", "", "Output file to write to (defaults to stdio)")
}
