package main

import (
	"log"
	"net/http/httputil"
	"net/url"
)

// NewProxy creates a reverse proxy for a target URL
func NewProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL %s: %v", target, err)
	}
	return httputil.NewSingleHostReverseProxy(targetURL)
}
