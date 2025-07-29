package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

var (
	clientsMu sync.Mutex                       // mutex to protect access to clients map
	clients   = make(map[chan string]struct{}) // map to track connected sse clients
)

const reloadScript = `
<script src="https://unpkg.com/morphdom"></script>
<script>
const evtSource = new EventSource("/events");
evtSource.onmessage = async function(event) {
    if (event.data === "reload") {
        const response = await fetch(window.location.href, { headers: { "X-Partial": "true" } });
        const parser = new DOMParser();
        const doc = parser.parseFromString(await response.text(), "text/html");
        morphdom(document.body, doc.body);
    }
};
</script>
`

// ssehandler sets up a server-sent events connection and streams messages to the client
func sseHandler(w http.ResponseWriter, r *http.Request) {
	// set headers for server-sent events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// ensure response supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	// create a message channel for this client
	msgCh := make(chan string)

	// add the client to the global clients map
	clientsMu.Lock()
	clients[msgCh] = struct{}{}
	clientsMu.Unlock()

	// remove client and clean up on disconnect
	defer func() {
		clientsMu.Lock()
		delete(clients, msgCh)
		clientsMu.Unlock()
		close(msgCh)
	}()

	// listen for messages and write them to the response
	for msg := range msgCh {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}
}

// broadcast sends a message to all connected clients
func broadcast(msg string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for ch := range clients {
		select {
		case ch <- msg: // try to send the message
		default: // if channel is blocked, close and remove it
			close(ch)
			delete(clients, ch)
		}
	}
}

// injectreloadscript wraps a handler and injects a reload script into html responses
func injectReloadScript(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// only handle html responses for injection
		if filepath.Ext(r.URL.Path) != ".html" && r.URL.Path != "/" {
			next.ServeHTTP(w, r)
			return
		}

		// capture the response body
		rw := &responseWriterCapture{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		// check if content type is html
		contentType := w.Header().Get("Content-Type")
		if contentType == "" || strings.Contains(contentType, "text/html") {
			html := rw.body.String()
			// inject the reload script before closing body tag
			modified := strings.Replace(html, "</body>", reloadScript+"</body>", 1)
			w.Header().Set("Content-Length", fmt.Sprint(len(modified)))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(modified))
		} else {
			// if not html, just pass the original response
			w.WriteHeader(rw.status)
			w.Write(rw.body.Bytes())
		}
	})
}

// responsewritercapture captures http responses for inspection and modification
type responseWriterCapture struct {
	http.ResponseWriter              // original response writer
	body                bytes.Buffer // buffer to store response body
	status              int          // status code to store
}

// write writes to the buffer instead of directly to the client
func (r *responseWriterCapture) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

// writeheader stores the status code instead of writing it immediately
func (r *responseWriterCapture) WriteHeader(statusCode int) {
	r.status = statusCode
}
