package link

import (
	"io"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	Href string
	Text string
}

type safeLinks struct {
	l  []Link
	mu sync.Mutex
}

func ParseLinks(r io.Reader) ([]Link, error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var l safeLinks
	var wg sync.WaitGroup
	wg.Add(1)
	go walkTreeRec(n, wg, &l)
	wg.Wait()

	return l.l, nil
}

func walkTreeRec(n *html.Node, wg sync.WaitGroup, l *safeLinks) {
	// After returning the links, signal to WaitGroup
	defer wg.Done()

	// Recursively call walkTreeRec on the child nodes
	for cn := n.FirstChild; cn != nil; cn = cn.NextSibling {
		wg.Add(1)
		go walkTreeRec(cn, wg, l)
	}

	// Check if current node is an <a> tag
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		var (
			link Link
			ok   bool
		)

		// Only if <a> tag has an `href` attribute can it be considered a link
		for _, attr := range n.Attr {
			if ok = attr.Key == "href"; ok {
				link.Href = attr.Val
				break
			}
		}

		// If it is a link, then append link to links
		if ok {
			l.mu.Lock()
			l.l = append(l.l, link)
			l.mu.Unlock()
		}
	}
}
