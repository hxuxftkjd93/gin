package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStreamDisconnect(t *testing.T) {
	exited := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &Context{
			Request: r,
			Writer:  &responseWriter{ResponseWriter: w},
		}

		c.Stream(func(w io.Writer) bool {
			select {
			case <-r.Context().Done():
				return false
			default:
				_, err := w.Write([]byte("data\n"))
				if err != nil {
					return false
				}
				time.Sleep(10 * time.Millisecond)
				return true
			}
		})
		close(exited)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 5)
	_, err = io.ReadFull(resp.Body, buf)
	if err != nil {
		t.Fatal(err)
	}

	cancel()

	select {
	case <-exited:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Stream did not exit after client disconnect")
	}
}