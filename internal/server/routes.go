package server

import (
	"net/http"

	"github.com/octyne/octyne/internal/types"
)

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", healthHandler)
	s.mux.Handle(
		"POST /v1/chat/completions",
		withRequestID(http.HandlerFunc(s.chatHandler)),
	)
	s.mux.Handle(
		"/v1/chat/completions",
		withRequestID(http.HandlerFunc(methodNotAllowedHandler)),
	)
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
