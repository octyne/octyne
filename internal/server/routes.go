package server

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", healthHandler)

}
