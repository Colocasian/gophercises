package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Colocasian/gophercises/urlshort"
	bolt "go.etcd.io/bbolt"
)

func main() {
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var handler http.Handler

	switch cli.format {
	case "bolt":
		db, err := bolt.Open(cli.conf, 0600, nil)
		if err != nil {
			fmt.Println("Oh no, could not open DB")
			panic(err)
		}
		defer db.Close()

		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("PathMap"))
			if err != nil {
				return fmt.Errorf("create bucket: %v", err)
			}
			b.Put([]byte("/urlshort"), []byte("https://github.com/gophercises/urlshort"))
			b.Put([]byte("/urlshort-final"), []byte("https://github.com/gophercises/urlshort/tree/solution"))
			return nil
		})
		if err != nil {
			panic(err)
		}

		handler = urlshort.BoltHandler(db, mapHandler)

	default:
		// Build the ConfigHandler using the mapHandler as the
		// fallback
		conf, err := os.ReadFile(cli.conf)
		if err != nil {
			panic(err)
		}

		configHandler, err := urlshort.ConfigHandler([]byte(conf), cli.format, mapHandler)
		if err != nil {
			panic(err)
		}
		handler = configHandler
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
