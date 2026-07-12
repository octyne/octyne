package server

import (
	"log"
	"net/http"
	"time"

	"github.com/octyne/octyne/internal/gateway"
	"github.com/octyne/octyne/internal/registry"
)

type Server struct {
	mux           *http.ServeMux
	gateway       *gateway.Service
	modelRegistry *registry.Registry
	httpServer    *http.Server
}

func New(addr string, gateway *gateway.Service, modelRegistry *registry.Registry) *Server {

	s := &Server{
		mux:           http.NewServeMux(),
		gateway:       gateway,
		modelRegistry: modelRegistry,
	}

	s.routes()

	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           s.mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      0, // Long-lived SSE streams must not have a global write deadline.
		IdleTimeout:       120 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	log.Printf("Octyne starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
