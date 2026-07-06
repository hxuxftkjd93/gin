package main

import (
	"io"
	"net/http"
)

// Context represents the context of the request.
type Context struct {
	Request *http.Request
	Writer  ResponseWriter
}

// Stream streams data.
func (c *Context) Stream(step func(w io.Writer) bool) bool {
	w := c.Writer
	for {
		select {
		case <-c.Request.Context().Done():
			return false
		default:
			if !step(w) {
				return true
			}
			w.Flush()
		}
	}
}