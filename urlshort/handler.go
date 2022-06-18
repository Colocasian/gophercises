package urlshort

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If request was not a GET request, fallback
		if r.Method != "GET" {
			fallback.ServeHTTP(w, r)
			return
		}

		// If redirect path not found in map, fallback
		url, ok := pathsToUrls[r.URL.Path]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}

		// Redirect to the correct URL
		http.Redirect(w, r, url, 302)
	}
}

type pathMap []struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func buildMap(pm pathMap) map[string]string {
	m := make(map[string]string)
	for _, r := range pm {
		m[r.Path] = r.Url
	}
	return m
}

// ConfigHandler will parse the provided config and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the config, then the
// fallback http.Handler will be called instead.
//
// Config is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid config data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func ConfigHandler(conf []byte, format string, fallback http.Handler) (http.HandlerFunc, error) {
	var pm pathMap

	var err error
	switch format {
	case "yaml":
		err = yaml.Unmarshal(conf, &pm)
	case "json":
		err = json.Unmarshal(conf, &pm)
	default:
		err = errors.New(fmt.Sprintf("unknown config format %q", format))
	}
	if err != nil {
		return nil, err
	}

	m := buildMap(pm)
	return MapHandler(m, fallback), nil
}
