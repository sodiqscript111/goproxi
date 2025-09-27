package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// Config structs
type Config struct {
	Proxy  ProxyConfig `mapstructure:"proxy"`
	Routes []Route     `mapstructure:"routes"`
}

type ProxyConfig struct {
	Bind string `mapstructure:"bind"`
}

type Route struct {
	Prefix string `mapstructure:"prefix"`
	Target string `mapstructure:"target"`
}

func main() {
	// Load config file
	viper.SetConfigName("goproxi") // file name without extension
	viper.SetConfigType("toml")
	viper.AddConfigPath(".") // look in current dir

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	// Build proxy handlers for each route
	handlers := make(map[string]http.Handler)
	for _, route := range config.Routes {
		handlers[route.Prefix] = NewProxy(route.Target)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		for prefix, proxy := range handlers {
			if strings.HasPrefix(r.URL.Path, prefix) {
				// Trim the prefix before forwarding
				r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
				if r.URL.Path == "" {
					r.URL.Path = "/"
				}
				proxy.ServeHTTP(w, r)
				return
			}
		}

		http.NotFound(w, r)
	})

	log.Printf("Starting proxy on :%s", config.Proxy.Bind)
	log.Fatal(http.ListenAndServe(":"+config.Proxy.Bind, nil))
}
