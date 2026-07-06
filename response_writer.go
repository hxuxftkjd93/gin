package main

import (
	"net/http"
)

// ResponseWriter wraps http.ResponseWriter.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	Status() int
	Size() int
	Written() bool
}

type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.status != 0
}

func (w *responseWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.WriteHeader(http.StatusOK)
	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

func (w *responseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}