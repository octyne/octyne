package server

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/octyne/octyne/internal/requestid"
)

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := newRequestID()
		if err != nil {
			writeOpenAIError(w, err)
			return
		}

		w.Header().Set("x-request-id", id)
		next.ServeHTTP(w, r.WithContext(requestid.WithContext(r.Context(), id)))
	})
}

func newRequestID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "req_" + hex.EncodeToString(bytes), nil
}
