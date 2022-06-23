package sitemap

import (
	"encoding/xml"
	"mime"
	"net/http"
	"net/url"
	"path"
	"sync"

	"github.com/Colocasian/gophercises/link"
)

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	NS      string   `xml:"xmlns,attr"`
	URLs    []URL
}

type URL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

func CrawlMap(ru string, d int) (*URLSet, error) {
	// Parse root URL
	root, err := url.Parse(ru)
	if err != nil {
		return nil, err
	}

	// Initialize the queues
	var cur, nxt []*url.URL
	cur = []*url.URL{root}
	var nmu sync.Mutex

	// Sitemap, containing the array of links
	var smarr []*url.URL
	var smu sync.Mutex

	var m sync.Map
	m.Store(root.String(), true)
	var wg sync.WaitGroup
	for i := 0; i < d && len(cur) != 0; i++ {
		wg.Add(len(cur))
		for _, u := range cur {
			go func(loc *url.URL) {
				// After completion, signal to WaitGroup
				defer wg.Done()

				resp, err := http.Get(loc.String())
				if err != nil {
					return
				}
				var closed bool
				defer func() {
					// Close the response body if not closed yet
					if !closed {
						resp.Body.Close()
					}
				}()

				// Check status code of response
				if resp.StatusCode < 200 || resp.StatusCode > 299 {
					return
				}

				// Append URL to sitemap
				smu.Lock()
				smarr = append(smarr, loc)
				smu.Unlock()

				// Check Content-Type header
				if !checkContentTypeHTML(resp) {
					return
				}

				// Extract links for response body, the close resp.Body
				links, err := link.ParseLinks(resp.Body)
				closed = true
				resp.Body.Close()
				if err != nil {
					return
				}

				for _, l := range links {
					ph, err := normalizeHref(loc, l.Href)
					if err != nil || ph.Host != root.Host {
						continue
					}

					if _, loaded := m.LoadOrStore(ph.String(), true); loaded {
						continue
					}
					nmu.Lock()
					nxt = append(nxt, ph)
					nmu.Unlock()
				}
			}(u)
		}
		wg.Wait()

		cur = nxt
		nxt = nil
	}

	var sm URLSet
	sm.NS = "http://www.sitemaps.org/schemas/sitemap/0.9"
	for _, u := range smarr {
		sm.URLs = append(sm.URLs, URL{Loc: u.String()})
	}
	return &sm, nil
}

func checkContentTypeHTML(resp *http.Response) bool {
	// Check Content-Type header
	// Common header variations for Content-Type
	headers := []string{"Content-Type", "content-type"}
	var (
		mimetype []string
		ok       bool
	)

	// Check for all the common variations
	for _, header := range headers {
		if mimetype, ok = resp.Header[header]; ok {
			break
		}
	}

	// If not found or multiple mimetypes on header, return false
	if !ok || len(mimetype) != 1 {
		return false
	}
	// Parse MIME to get media type, then check it for text/html (no MIME sniffing)
	mediatype, _, err := mime.ParseMediaType(mimetype[0])
	if err != nil || mediatype != "text/html" {
		return false
	}

	return true
}

func normalizeHref(from *url.URL, href string) (*url.URL, error) {
	ph, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	// Add domain name if missing
	if ph.Scheme == "" && ph.User == nil && ph.Host == "" {
		ph.Scheme = from.Scheme
		ph.Host = from.Host
		if len(href) == 0 || href[0] == '/' {
			ph.Path = href
		} else {
			ph.Path = path.Join("/", from.Path, ph.Path)
		}
	}

	return ph, nil
}
