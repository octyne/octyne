package server

import (
	"log"
	"net/http"

	"github.com/octyne/octyne/internal/gateway"
	"github.com/octyne/octyne/internal/registry"
)

type Server struct {
	mux           *http.ServeMux
	gateway       *gateway.Service
	modelRegistry *registry.Registry
}

func New(gateway *gateway.Service, modelRegistry *registry.Registry) *Server {

	s := &Server{
		mux:           http.NewServeMux(),
		gateway:       gateway,
		modelRegistry: modelRegistry,
	}

	s.routes()

	return s
}

func (s *Server) Start(addr string) error {
	log.Printf("Octyne starting on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
