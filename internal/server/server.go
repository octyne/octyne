package server

import (
	"log"
	"net/http"
	"time"

	"github.com/usekeel/keel/internal/gateway"
	"github.com/usekeel/keel/internal/providers"
)

type Server struct {
	mux     *http.ServeMux
	gateway *gateway.Service
}

func New() *Server {
	providerRegistry := providers.NewRegistry()

	providerRegistry.Register(
		"openai",
		&providers.Provider{
			Name: "openai",
			Config: providers.Config{
				Name:    "openai",
				BaseURL: "https://api.openai.com/v1",
				Timeout: 30 * time.Second,
			},
		},
	)

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
