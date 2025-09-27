package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy creates a reverse proxy for a target URL
func NewProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL %s: %v", target, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Save the original director
	originalDirector := proxy.Director

	// Override the director with custom logic
	proxy.Director = func(req *http.Request) {
		// Call the original director first
		originalDirector(req)

		// Add your custom header
		req.Header.Add("X-My-Proxi", "goproxi")
	}

	return proxy
}
