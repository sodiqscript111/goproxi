package main

import (
	"log"
	"net/http/httputil"
	"net/url"
)

func NewProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL %s: %v", target, err)
	}
	return httputil.NewSingleHostReverseProxy(targetURL)
}
