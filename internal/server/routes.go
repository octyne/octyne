package server

import (
	"net/http"

	"github.com/octyne/octyne/internal/types"
)

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", healthHandler)
	s.mux.HandleFunc("GET /v1/models", s.modelsHandler)
	s.mux.HandleFunc("/v1/models", methodNotAllowedHandler)
	s.mux.HandleFunc("POST /v1/chat/completions", s.chatHandler)
	s.mux.HandleFunc("/v1/chat/completions", methodNotAllowedHandler)
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	code := "method_not_allowed"
	writeOpenAIError(w, &types.APIError{
		Kind:       types.ErrorKindInvalidRequest,
		Message:    "Method not allowed.",
		Code:       &code,
		HTTPStatus: http.StatusMethodNotAllowed,
	})
}
