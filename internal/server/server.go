package server

import (
	"log"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func New() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}

	s.routes()

	return s
}

func (s *Server) Start(addr string) error {
	log.Printf("Keel starting on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
