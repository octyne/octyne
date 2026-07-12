package server

import (
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
