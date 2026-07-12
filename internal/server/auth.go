package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/octyne/octyne/internal/auth"
	"github.com/octyne/octyne/internal/types"
)

func withAuthentication(
	verifier auth.Verifier,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if verifier == nil {
			writeOpenAIError(
				w,
				errors.New("authentication verifier is not configured"),
			)
			return
		}

		key, ok := bearerKey(r)
		if !ok {
			writeAuthenticationError(w)
			return
		}

		valid, err := verifier.Verify(r.Context(), key)
		if err != nil {
			writeOpenAIError(w, err)
			return
		}
		if !valid {
			writeAuthenticationError(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bearerKey(r *http.Request) (string, bool) {
	values := r.Header.Values("Authorization")
	if len(values) != 1 {
		return "", false
	}

	scheme, key, found := strings.Cut(values[0], " ")
	if !found ||
		!strings.EqualFold(scheme, "Bearer") ||
		key == "" ||
		strings.TrimSpace(key) != key ||
		strings.ContainsAny(key, " \t\r\n") {
		return "", false
	}

	return key, true
}

func writeAuthenticationError(w http.ResponseWriter) {
	code := "invalid_api_key"

	w.Header().Set("WWW-Authenticate", `Bearer realm="octyne"`)
	writeOpenAIError(w, &types.APIError{
		Kind:       types.ErrorKindAuthentication,
		Message:    "Invalid API key.",
		Code:       &code,
		HTTPStatus: http.StatusUnauthorized,
	})
}
