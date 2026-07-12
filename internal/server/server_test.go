package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewConfiguresHTTPServer(t *testing.T) {
	server := New(":4321", nil, nil)

	if server.httpServer.Addr != ":4321" {
		t.Errorf("Addr = %q, want %q", server.httpServer.Addr, ":4321")
	}
	if server.httpServer.Handler != server.mux {
		t.Error("Handler does not use the server mux")
	}

	tests := []struct {
		name string
		got  time.Duration
		want time.Duration
	}{
		{name: "read header timeout", got: server.httpServer.ReadHeaderTimeout, want: 5 * time.Second},
		{name: "read timeout", got: server.httpServer.ReadTimeout, want: 30 * time.Second},
		{name: "write timeout", got: server.httpServer.WriteTimeout, want: 0},
		{name: "idle timeout", got: server.httpServer.IdleTimeout, want: 120 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("timeout = %s, want %s", tt.got, tt.want)
			}
		})
	}
}

func TestRunReturnsServeError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := New("invalid address", nil, nil).Run(ctx)
	if err == nil {
		t.Fatal("Run error = nil, want listen error")
	}
	if !strings.Contains(err.Error(), "serve HTTP server") {
		t.Errorf("Run error = %q, want serve operation context", err)
	}
}

func TestRunDrainsActiveRequestOnCancellation(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve server address: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("close reserved listener: %v", err)
	}

	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})
	released := false
	defer func() {
		if !released {
			close(releaseRequest)
		}
	}()

	server := New(addr, nil, nil)
	server.httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(requestStarted)
		<-releaseRequest
		w.WriteHeader(http.StatusNoContent)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runErr := make(chan error, 1)
	go func() {
		runErr <- server.Run(ctx)
	}()

	client := &http.Client{Timeout: 2 * time.Second}
	requestErr := make(chan error, 1)
	go func() {
		deadline := time.Now().Add(2 * time.Second)
		for {
			resp, err := client.Get("http://" + addr)
			if err == nil {
				_ = resp.Body.Close()
				requestErr <- nil
				return
			}
			if time.Now().After(deadline) {
				requestErr <- fmt.Errorf("send request: %w", err)
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	select {
	case <-requestStarted:
	case err := <-requestErr:
		t.Fatalf("request completed before handler started: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("active request did not reach handler")
	}

	cancel()

	select {
	case err := <-runErr:
		t.Fatalf("Run returned before active request completed: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	close(releaseRequest)
	released = true

	if err := <-requestErr; err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-runErr:
		if err != nil {
			t.Fatalf("Run error = %v, want nil", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after active request completed")
	}
}
