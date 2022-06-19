package story

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (s *Story) HTTPHandler(rt, at *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}

		path := filepath.Clean(r.URL.Path)

		switch path {
		case "/":
			rt.Execute(w, s)
		default:
			arc, ok := (*s)[r.URL.Path[1:]]
			if !ok {
				http.NotFound(w, r)
				return
			}
			at.Execute(w, arc)
		}
	}
}
