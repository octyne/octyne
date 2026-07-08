package server

import (
	"log"
	"net/http"

	"github.com/octyne/octyne/internal/gateway"
)

type Server struct {
	mux     *http.ServeMux
	gateway *gateway.Service
}

func New(gateway *gateway.Service) *Server {

	s := &Server{
		mux:     http.NewServeMux(),
		gateway: gateway,
	}

	s.routes()

	return s
}

func (s *Server) Start(addr string) error {
	log.Printf("Octyne starting on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
