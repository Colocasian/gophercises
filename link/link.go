package link

import (
	"io"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	Href string
	Text string
}

func ParseLinks(r io.Reader) ([]Link, error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var l safeLinks
	var wg sync.WaitGroup
	wg.Add(1)
	go walkTreeRec(n, &wg, &l)
	wg.Wait()

	return l.l, nil
}

type safeLinks struct {
	l  []Link
	mu sync.Mutex
}

func walkTreeRec(n *html.Node, wg *sync.WaitGroup, l *safeLinks) {
	// After returning the links, signal to WaitGroup
	defer wg.Done()

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
				// Sequential code to get link text
				var t strings.Builder
				getDescendantText(n, &t)
				link.Text = t.String()
				// Code to get link text concurrently:
				//var b []byte
				//getDescendantTextConc(n, &b, nil)
				//link.Text = string(b)
				break
			}
		}

		// If it is a link, then append link to links
		if ok {
			l.mu.Lock()
			l.l = append(l.l, link)
			l.mu.Unlock()
		}

		// Do not process links inside <a> tags: invalid HTML
		return
	}

	// Recursively call walkTreeRec on the child nodes
	for cn := n.FirstChild; cn != nil; cn = cn.NextSibling {
		wg.Add(1)
		go walkTreeRec(cn, wg, l)
	}
}

func getDescendantText(n *html.Node, t *strings.Builder) {
	if n.Type == html.TextNode {
		t.WriteString(n.Data)
	}

	// Recursively call getDescendantText on the child nodes
	for cn := n.FirstChild; cn != nil; cn = cn.NextSibling {
		getDescendantText(cn, t)
	}
}

func getDescendantTextConc(n *html.Node, t *strings.Builder, wg *sync.WaitGroup) {
	// If wg is not null, after returning the links, signal to WaitGroup
	if wg != nil {
		defer wg.Done()
	}

	var sz int
	for cn := n.FirstChild; cn != nil; cn = cn.NextSibling {
		sz++
	}

	cts := make([]strings.Builder, sz)
	var cwg sync.WaitGroup
	// Recursively call getDescendantText on the child nodes
	for i, cn := 0, n.FirstChild; cn != nil; i, cn = i+1, cn.NextSibling {
		cwg.Add(1)
		go getDescendantTextConc(cn, &cts[i], &cwg)
	}
	// Wait for all the children to return text
	cwg.Wait()

	if n.Type == html.TextNode {
		t.WriteString(n.Data)
	}
	for _, ct := range cts {
		t.WriteString(ct.String())
	}
}
