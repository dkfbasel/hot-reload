package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// reverseproxyhandler creates a reverse proxy to forward requests to the target server
func reverseProxyHandler(target string) http.Handler {
	// parse the target url
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("error parsing target url: %v", err)
	}

	// create a reverse proxy for the target
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	return proxy
}

// runhttpserver sets up the http server, sse handler, and proxy routing
func runHttpServer() {
	// target server to proxy requests to
	target := "http://localhost:8081"

	// create the reverse proxy handler
	proxy := reverseProxyHandler(target)

	// handle sse hotreload endpoint
	http.HandleFunc("/hotreload", sseHandler)

	// wrap proxy with script injector and set as root handler
	http.Handle("/", injectReloadScript(proxy))

	// start the server on port 3333
	log.Println("server started at http://localhost:3333")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
