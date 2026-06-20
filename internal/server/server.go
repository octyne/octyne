package server

import (
	"log"
	"net/http"

	"github.com/usekeel/keel/internal/gateway"
	"github.com/usekeel/keel/internal/providers"
)

type Server struct {
	mux     *http.ServeMux
	gateway *gateway.Service
}

func New() *Server {
	providerRegistry := providers.NewRegistry()

	gateway := gateway.New(providerRegistry)

	s := &Server{
		mux:     http.NewServeMux(),
		gateway: gateway,
	}

	s.routes()

	return s
}

func (s *Server) Start(addr string) error {
	log.Printf("Keel starting on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
