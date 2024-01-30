package util

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type HijackableResponseWriter struct {
	http.ResponseWriter
	hijacked bool
	status   int
}

func (h *HijackableResponseWriter) WriteHeader(statusCode int) {
	if !h.hijacked {
		h.ResponseWriter.WriteHeader(statusCode)
		h.status = statusCode
	}
}

func (h *HijackableResponseWriter) Status() int {
	return h.status
}

func (h *HijackableResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.hijacked = true
	hj, ok := h.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("webserver doesn't support hijacking")
	}
	return hj.Hijack()
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rr := &HijackableResponseWriter{ResponseWriter: w}
		next.ServeHTTP(rr, r)
		status := rr.Status()
		if !rr.hijacked {
			fmt.Printf("Received request: %s %s - Status: %d\n", r.Method, r.RequestURI, status)
		}
	})
}
