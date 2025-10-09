package server

import (
	"context"
	"goproxi/internals/config"
	"goproxi/internals/proxy"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func Start() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("There was an error", err.Error())
	}

	// Build proxy handlers for each route
	handlers := make(map[string]http.Handler)
	for _, route := range cfg.Routes {
		handlers[route.Prefix] = proxy.NewProxy(route.Target)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		for prefix, p := range handlers {
			if strings.HasPrefix(r.URL.Path, prefix) {
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

	srv := &http.Server{
		Addr:         ":" + cfg.Proxy.Bind,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting proxy on :%s", cfg.Proxy.Bind)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
