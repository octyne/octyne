package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/octyne/octyne/internal/auth"
)

func TestAuthenticationMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		authorization []string
		query         string
		wantStatus    int
		wantNext      bool
	}{
		{
			name:       "missing authorization",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:          "wrong scheme",
			authorization: []string{"Basic " + testClientAPIKey},
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:          "missing bearer token",
			authorization: []string{"Bearer"},
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:          "empty bearer token",
			authorization: []string{"Bearer "},
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:          "incorrect bearer token",
			authorization: []string{"Bearer incorrect-client-key-with-at-least-32-characters"},
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:          "multiple authorization headers",
			authorization: []string{"Bearer " + testClientAPIKey, "Bearer " + testClientAPIKey},
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:       "query parameter is not accepted",
			query:      "?api_key=" + testClientAPIKey,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:          "extra separator space",
			authorization: []string{"Bearer  " + testClientAPIKey},
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:          "valid bearer token",
			authorization: []string{"Bearer " + testClientAPIKey},
			wantStatus:    http.StatusNoContent,
			wantNext:      true,
		},
		{
			name:          "case insensitive bearer scheme",
			authorization: []string{"bearer " + testClientAPIKey},
			wantStatus:    http.StatusNoContent,
			wantNext:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextCalled := false
			handler := withAuthentication(
				newTestVerifier(),
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					nextCalled = true
					w.WriteHeader(http.StatusNoContent)
				}),
			)
			request := httptest.NewRequest(http.MethodGet, "/v1/models"+test.query, nil)
			for _, value := range test.authorization {
				request.Header.Add("Authorization", value)
			}
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
			if nextCalled != test.wantNext {
				t.Errorf("next called = %t, want %t", nextCalled, test.wantNext)
			}
			if test.wantStatus == http.StatusUnauthorized {
				assertAuthenticationError(t, recorder)
			}
		})
	}
}

func TestAuthenticationMiddlewareFailsClosed(t *testing.T) {
	tests := []struct {
		name     string
		verifier auth.Verifier
	}{
		{
			name:     "missing verifier",
			verifier: nil,
		},
		{
			name:     "verifier error",
			verifier: errorVerifier{err: errors.New("credential store unavailable")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextCalled := false
			handler := withAuthentication(
				test.verifier,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					nextCalled = true
				}),
			)
			request := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
			authenticateRequest(request)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
			}
			if nextCalled {
				t.Error("next handler was called after authentication failure")
			}
			if strings.Contains(recorder.Body.String(), "credential store unavailable") {
				t.Error("response exposes verifier error details")
			}
		})
	}
}

func TestAuthenticationRouteBoundary(t *testing.T) {
	server := New(":0", newTestLogger(), nil, nil, newTestVerifier())

	tests := []struct {
		name          string
		path          string
		authenticated bool
		wantStatus    int
	}{
		{
			name:       "health is public",
			path:       "/health",
			wantStatus: http.StatusOK,
		},
		{
			name:       "models requires authentication",
			path:       "/v1/models",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "exact v1 root requires authentication",
			path:       "/v1",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "unknown v1 route requires authentication",
			path:       "/v1/unknown",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:          "authenticated unknown v1 route reaches router",
			path:          "/v1/unknown",
			authenticated: true,
			wantStatus:    http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			if test.authenticated {
				authenticateRequest(request)
			}
			recorder := httptest.NewRecorder()

			server.httpServer.Handler.ServeHTTP(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
			if got := recorder.Header().Get("x-request-id"); !strings.HasPrefix(got, "req_") {
				t.Errorf("x-request-id = %q, want req_ prefix", got)
			}
		})
	}
}

func TestAuthenticationDoesNotForwardClientKeyUpstream(t *testing.T) {
	var upstreamAuthorization string
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamAuthorization = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"chatcmpl-auth","object":"chat.completion","created":1,"model":"gpt-5-nano","choices":[]}`)
	}))
	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(`{"model":"openai/gpt-5-nano","messages":[{"role":"user","content":"hello"}]}`),
	)
	authenticateRequest(request)
	recorder := httptest.NewRecorder()

	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if upstreamAuthorization != "Bearer test-provider-api-key" {
		t.Errorf("upstream Authorization = %q, want provider credential", upstreamAuthorization)
	}
	if strings.Contains(upstreamAuthorization, testClientAPIKey) {
		t.Error("upstream Authorization contains the Octyne client key")
	}
}

func TestAuthenticationLoggingDoesNotExposeClientKey(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	server := New(":0", logger, nil, nil, newTestVerifier())
	request := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	request.Header.Set("Authorization", "Bearer incorrect-client-key-with-at-least-32-characters")
	recorder := httptest.NewRecorder()

	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
	if strings.Contains(output.String(), "incorrect-client-key") {
		t.Error("request log exposes the presented client key")
	}
}

func assertAuthenticationError(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	if got := recorder.Header().Get("WWW-Authenticate"); got != `Bearer realm="octyne"` {
		t.Errorf("WWW-Authenticate = %q, want Bearer challenge", got)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	var envelope openAIErrorEnvelope
	if err := json.NewDecoder(recorder.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode authentication error: %v", err)
	}
	if envelope.Error.Type != "authentication_error" ||
		envelope.Error.Message != "Invalid API key." ||
		value(envelope.Error.Code) != "invalid_api_key" {
		t.Errorf("unexpected authentication error: %+v", envelope.Error)
	}
}

type errorVerifier struct {
	err error
}

func (v errorVerifier) Verify(context.Context, string) (bool, error) {
	return false, v.err
}
