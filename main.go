package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	httpbinProxy := NewProxy("https://httpbin.org")
	golangProxy := NewProxy("https://golang.org")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		switch {
		case strings.HasPrefix(r.URL.Path, "/httpbin"):
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/httpbin")
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			httpbinProxy.ServeHTTP(w, r)

		case strings.HasPrefix(r.URL.Path, "/golang"):
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/golang")
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			golangProxy.ServeHTTP(w, r)

		default:
			http.NotFound(w, r)
		}
	})

	log.Println("Starting proxy on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
