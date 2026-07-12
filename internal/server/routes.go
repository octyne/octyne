package server

import (
	"net/http"

	"github.com/octyne/octyne/internal/auth"
	"github.com/octyne/octyne/internal/types"
)

func (s *Server) routes(verifier auth.Verifier) {
	s.mux.HandleFunc("GET /health", healthHandler)

	v1Mux := http.NewServeMux()
	v1Mux.HandleFunc("GET /v1/models", s.modelsHandler)
	v1Mux.HandleFunc("/v1/models", methodNotAllowedHandler)
	v1Mux.HandleFunc("POST /v1/chat/completions", s.chatHandler)
	v1Mux.HandleFunc("/v1/chat/completions", methodNotAllowedHandler)

	protectedV1 := withAuthentication(verifier, v1Mux)
	s.mux.Handle("/v1", protectedV1)
	s.mux.Handle("/v1/", protectedV1)
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
