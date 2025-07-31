package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	clientsMu sync.Mutex                       // mutex to protect access to clients map
	clients   = make(map[chan string]struct{}) // map to track connected sse clients
)

const reloadScript = `
<script src="https://unpkg.com/morphdom@2.6.1/dist/morphdom-umd.min.js"></script>
<script>
const evtSource = new EventSource("/hotreload");
evtSource.onmessage = async function(event) {
    if (event.data === "reload") {
		console.log("Reloading body dom...");
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

		if strings.HasPrefix(r.URL.Path, "/hotreload") || r.Header.Get("Accept") == "text/event-stream" {
			log.Println("Skipping script injection for SSE or hotreload endpoint")
			next.ServeHTTP(w, r)
			return
		}
		// capture the response
		rw := &responseWriterCapture{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		// check if the content type is html
		contentType := rw.Header().Get("Content-Type")
		if contentType == "" || strings.Contains(contentType, "text/html") {
			html := rw.body.String()
			// inject the reload script before </head>
			modified := strings.Replace(html, "</head>", reloadScript+"</head>", 1)
			w.Header().Set("Content-Length", fmt.Sprint(len(modified)))
			if !rw.wroteHeader {
				w.WriteHeader(http.StatusOK) // or rw.status if you want the original
			}
			w.Write([]byte(modified))
		} else {
			// just forward the original response
			w.WriteHeader(rw.status)
			w.Write(rw.body.Bytes())
		}
	})
}

// responsewritercapture captures http responses for inspection and modification
type responseWriterCapture struct {
	http.ResponseWriter              // original response writer
	body                bytes.Buffer // buffer to store response body
	status              int          // captured status code
	wroteHeader         bool         // tracks if writeheader was called
}

// write captures the response body (but does not write to client yet)
func (r *responseWriterCapture) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

// writeheader captures the status code
func (r *responseWriterCapture) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}
	r.status = statusCode
	r.wroteHeader = true
}

// flushtoclient sends the captured body to the real responsewriter
func (r *responseWriterCapture) FlushToClient() error {
	// if no status was written, default to 200
	if !r.wroteHeader {
		r.status = http.StatusOK
		r.ResponseWriter.WriteHeader(r.status)
	}
	_, err := r.ResponseWriter.Write(r.body.Bytes())
	return err
}

// implements http.hijacker
func (r *responseWriterCapture) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("underlying responsewriter does not support hijacking")
	}
	return hj.Hijack()
}

// implements http.flusher
func (r *responseWriterCapture) Flush() {
	if fl, ok := r.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

// implements http.closenotifier (deprecated but still seen)
func (r *responseWriterCapture) CloseNotify() <-chan bool {
	if cn, ok := r.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return make(chan bool)
}
