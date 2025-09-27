package server

import (
	"goproxi/internals/config"
	"goproxi/internals/proxy"
	"log"
	"net/http"
	"strings"
)

func Start() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("There was an error", err.Error())
	}

	// Build proxy handlers for each route
	handlers := make(map[string]http.Handler)
	for _, route := range cfg.Routes {
		handlers[route.Prefix] = proxy.NewProxy(route.Target)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		for prefix, p := range handlers {
			if strings.HasPrefix(r.URL.Path, prefix) {
				// Trim prefix before forwarding
				r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
				if r.URL.Path == "" {
					r.URL.Path = "/"
				}
				p.ServeHTTP(w, r)
				return
			}
		}

		http.NotFound(w, r)
	})

	log.Printf("Starting proxy on :%s", cfg.Proxy.Bind)
	log.Fatal(http.ListenAndServe(":"+cfg.Proxy.Bind, nil))
}
