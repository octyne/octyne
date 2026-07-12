package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/octyne/octyne/internal/types"
)

func TestOpenAIErrorTypeMapping(t *testing.T) {
	tests := []struct {
		kind types.ErrorKind
		want string
	}{
		{types.ErrorKindInvalidRequest, "invalid_request_error"},
		{types.ErrorKindNotFound, "invalid_request_error"},
		{types.ErrorKindAuthentication, "authentication_error"},
		{types.ErrorKindPermission, "permission_error"},
		{types.ErrorKindRateLimit, "rate_limit_error"},
		{types.ErrorKindTimeout, "server_error"},
		{types.ErrorKindUnavailable, "server_error"},
		{types.ErrorKindInternal, "server_error"},
	}

	for _, tt := range tests {
		t.Run(string(tt.kind), func(t *testing.T) {
			got := toOpenAIErrorEnvelope(&types.APIError{
				Kind:    tt.kind,
				Message: "message",
			})
			if got.Error.Type != tt.want {
				t.Errorf("type = %q, want %q", got.Error.Type, tt.want)
			}
		})
	}
}

func TestChatHandlerReturnsOpenAIErrorsAndRequestIDs(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantType   string
		wantParam  string
		wantCode   string
	}{
		{
			name:       "invalid JSON",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantType:   "invalid_request_error",
			wantCode:   "invalid_json",
		},
		{
			name:       "trailing JSON value",
			body:       `{"model":"openai/gpt-5-nano","messages":[{"role":"user","content":"hello"}]} {}`,
			wantStatus: http.StatusBadRequest,
			wantType:   "invalid_request_error",
			wantCode:   "invalid_json",
		},
		{
			name:       "missing model",
			body:       `{"messages":[{"role":"user","content":"hello"}]}`,
			wantStatus: http.StatusBadRequest,
			wantType:   "invalid_request_error",
			wantParam:  "model",
		},
		{
			name:       "missing messages",
			body:       `{"model":"openai/gpt-5-nano"}`,
			wantStatus: http.StatusBadRequest,
			wantType:   "invalid_request_error",
			wantParam:  "messages",
		},
		{
			name:       "unknown model",
			body:       `{"model":"not-registered","messages":[{"role":"user","content":"hello"}]}`,
			wantStatus: http.StatusNotFound,
			wantType:   "invalid_request_error",
			wantParam:  "model",
			wantCode:   "model_not_found",
		},
		{
			name:       "unqualified model",
			body:       `{"model":"gpt-5-nano","messages":[{"role":"user","content":"hello"}]}`,
			wantStatus: http.StatusNotFound,
			wantType:   "invalid_request_error",
			wantParam:  "model",
			wantCode:   "model_not_found",
		},
		{
			name:       "wrong HTTP method",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
			wantType:   "invalid_request_error",
			wantCode:   "method_not_allowed",
		},
	}

	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("upstream was called for a rejected request")
	}))

	seenIDs := make(map[string]bool)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := tt.method
			if method == "" {
				method = http.MethodPost
			}
			request := httptest.NewRequest(
				method,
				"/v1/chat/completions",
				strings.NewReader(tt.body),
			)
			authenticateRequest(request)
			recorder := httptest.NewRecorder()

			server.httpServer.Handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
			requestID := recorder.Header().Get("x-request-id")
			if !strings.HasPrefix(requestID, "req_") {
				t.Errorf("x-request-id = %q, want req_ prefix", requestID)
			}
			if seenIDs[requestID] {
				t.Errorf("x-request-id %q was reused", requestID)
			}
			seenIDs[requestID] = true

			var envelope openAIErrorEnvelope
			if err := json.NewDecoder(recorder.Body).Decode(&envelope); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if envelope.Error.Type != tt.wantType {
				t.Errorf("type = %q, want %q", envelope.Error.Type, tt.wantType)
			}
			if value(envelope.Error.Param) != tt.wantParam {
				t.Errorf("param = %q, want %q", value(envelope.Error.Param), tt.wantParam)
			}
			if value(envelope.Error.Code) != tt.wantCode {
				t.Errorf("code = %q, want %q", value(envelope.Error.Code), tt.wantCode)
			}
		})
	}
}

func TestChatHandlerMapsUpstreamOpenAIError(t *testing.T) {
	var clientRequestID string
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientRequestID = r.Header.Get("X-Client-Request-Id")
		w.Header().Set("x-request-id", "req_provider")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"error":{"message":"Too many requests.","type":"rate_limit_error","param":null,"code":"rate_limit_exceeded"}}`)
	}))

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(`{"model":"openai/gpt-5-nano","messages":[{"role":"user","content":"hello"}]}`),
	)
	authenticateRequest(request)
	recorder := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", recorder.Code)
	}
	octyneRequestID := recorder.Header().Get("x-request-id")
	if octyneRequestID == "" || clientRequestID != octyneRequestID {
		t.Errorf("upstream X-Client-Request-Id = %q, want %q", clientRequestID, octyneRequestID)
	}
	if octyneRequestID == "req_provider" {
		t.Error("provider request ID replaced Octyne request ID")
	}

	var envelope openAIErrorEnvelope
	if err := json.NewDecoder(recorder.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if envelope.Error.Type != "rate_limit_error" ||
		envelope.Error.Message != "Too many requests." ||
		value(envelope.Error.Code) != "rate_limit_exceeded" {
		t.Errorf("unexpected error envelope: %+v", envelope.Error)
	}
}

func TestChatHandlerStreamsOpenAIErrorEvent(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"error\":{\"message\":\"stream failed\",\"type\":\"server_error\",\"param\":null,\"code\":\"stream_error\"}}\n\n")
	}))

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(`{"model":"openai/gpt-5-nano","messages":[{"role":"user","content":"hello"}],"stream":true}`),
	)
	authenticateRequest(request)
	recorder := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	want := `data: {"error":{"message":"stream failed","type":"server_error","param":null,"code":"stream_error"}}` + "\n\n"
	if recorder.Body.String() != want {
		t.Errorf("body = %q, want %q", recorder.Body.String(), want)
	}
	if strings.Contains(recorder.Body.String(), "[DONE]") {
		t.Error("failed stream emitted [DONE]")
	}
}

func value(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
