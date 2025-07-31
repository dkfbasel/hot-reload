package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
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

	// return a handler that retries on error
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts := 0
		// retry the proxy if the target server is temporarily unavailable
		for {
			if !checkTargetHealth(target) {
				// if the target server is unavailable, retry a few times
				if attempts >= 5 {
					// after 5 retries, give up and return service unavailable
					http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
					return
				}

				time.Sleep(2 * time.Second)
				attempts++
				continue
			}

			// if the server is up, proxy the request
			proxy.ServeHTTP(w, r)
			return
		}
	})
}

// runhttpserver sets up the http server, sse handler, and proxy routing
func runHttpServer(proxyTarget string, proxyPort string) {

	// create the reverse proxy handler pointing to the target address
	proxy := reverseProxyHandler(proxyTarget)

	// handle sse hotreload endpoint
	http.HandleFunc("/hotreload", sseHandler)

	// wrap proxy with script injector and set as root handler
	http.Handle("/", injectReloadScript(proxy))

	// start the server on port 3333
	log.Println("server started at http://localhost:" + proxyPort)
	log.Fatal(http.ListenAndServe(":"+proxyPort, nil))
}

// checkTargetHealth checks if the target server is up and reachable.
func checkTargetHealth(target string) bool {
	resp, err := http.Get(target)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()
	return true
}
