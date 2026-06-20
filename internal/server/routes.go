package server

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", healthHandler)
	s.mux.HandleFunc("POST /v1/chat/completions", s.chatHandler)
}
