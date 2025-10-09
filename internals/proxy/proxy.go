package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL %s: %v", target, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {

		originalDirector(req)

		req.Header.Add("X-My-Proxi", "goproxi")
	}

	return proxy
}
