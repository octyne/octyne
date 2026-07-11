package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/octyne/octyne/internal/adapters/openai"
	"github.com/octyne/octyne/internal/gateway"
	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/types"
)

func TestChatHandlerStreamsOpenAICompatibleSSE(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/chat/completions" {
				t.Errorf(
					"path = %q, want /chat/completions",
					r.URL.Path,
				)
			}

			if got := r.Header.Get("Accept"); got != "text/event-stream" {
				t.Errorf(
					"Accept = %q, want text/event-stream",
					got,
				)
			}

			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read request body: %v", err)
				return
			}

			if !strings.Contains(
				string(requestBody),
				`"stream":true`,
			) {
				t.Errorf(
					"request body does not enable streaming: %s",
					requestBody,
				)
			}

			w.Header().Set("Content-Type", "text/event-stream")

			_, err = io.WriteString(w, `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":123,"model":"gpt-5-nano","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null,"logprobs":null}]}

data: [DONE]

`)
			if err != nil {
				t.Errorf("write upstream stream: %v", err)
			}
		},
	))

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(
			`{"model":"gpt-5-nano","messages":[{"role":"user","content":"Hello"}],"stream":true}`,
		),
	)

	recorder := httptest.NewRecorder()
	server.mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d",
			recorder.Code,
			http.StatusOK,
		)
	}

	if got := recorder.Header().Get("Content-Type"); got != "text/event-stream" {
		t.Errorf(
			"Content-Type = %q, want text/event-stream",
			got,
		)
	}

	wantBody := `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":123,"model":"gpt-5-nano","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null,"logprobs":null}]}

data: [DONE]

`

	if got := recorder.Body.String(); got != wantBody {
		t.Errorf(
			"response body:\n%s\nwant:\n%s",
			got,
			wantBody,
		)
	}

	if !recorder.Flushed {
		t.Error("response was not flushed")
	}
}

func TestChatHandlerReturnsOpenAICompatibleJSON(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read request body: %v", err)
				return
			}

			if strings.Contains(string(requestBody), `"stream":true`) {
				t.Errorf(
					"non-streaming request enables streaming: %s",
					requestBody,
				)
			}

			w.Header().Set("Content-Type", "application/json")

			_, err = io.WriteString(w, `{
				"id":"chatcmpl-456",
				"object":"chat.completion",
				"created":456,
				"model":"gpt-5-nano",
				"choices":[{
					"index":0,
					"message":{"role":"assistant","content":"Hello"},
					"finish_reason":"stop",
					"logprobs":null
				}],
				"usage":{"prompt_tokens":5,"completion_tokens":1,"total_tokens":6}
			}`)
			if err != nil {
				t.Errorf("write upstream response: %v", err)
			}
		},
	))

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(
			`{"model":"gpt-5-nano","messages":[{"role":"user","content":"Hello"}]}`,
		),
	)

	recorder := httptest.NewRecorder()
	server.mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d: %s",
			recorder.Code,
			http.StatusOK,
			recorder.Body.String(),
		)
	}

	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	var response types.ChatCompletionResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.ID != "chatcmpl-456" ||
		response.Object != "chat.completion" ||
		response.Created != 456 ||
		response.Model != "gpt-5-nano" {
		t.Errorf("unexpected response metadata: %+v", response)
	}

	if len(response.Choices) != 1 {
		t.Fatalf("len(Choices) = %d, want 1", len(response.Choices))
	}

	choice := response.Choices[0]
	if choice.Index != 0 ||
		choice.Message.Role != "assistant" ||
		choice.Message.Content != "Hello" ||
		choice.FinishReason == nil ||
		*choice.FinishReason != "stop" {
		t.Errorf("unexpected choice: %+v", choice)
	}

	if response.Usage == nil {
		t.Error("Usage = nil, want OpenAI usage object")
	}
}

func newTestServer(t *testing.T, upstreamHandler http.Handler) *Server {
	t.Helper()

	upstream := httptest.NewServer(upstreamHandler)
	t.Cleanup(upstream.Close)

	config := providers.Config{
		Name:                           "openai",
		BaseURL:                        upstream.URL,
		NonStreamingTimeout:            time.Second,
		StreamingResponseHeaderTimeout: time.Second,
	}

	providerRegistry := providers.NewRegistry()
	providerRegistry.Register(
		"openai",
		providers.New(
			config,
			openai.New(config),
		),
	)

	return New(gateway.New(providerRegistry))
}

func TestChatHandlerRejectsEmptyMessages(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t.Error("upstream was called for an invalid request")
			http.Error(w, "unexpected upstream call", http.StatusInternalServerError)
		},
	))

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(
			`{"model":"gpt-5-nano","messages":[]}`,
		),
	)

	recorder := httptest.NewRecorder()
	server.mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(
			"status = %d, want %d",
			recorder.Code,
			http.StatusBadRequest,
		)
	}

	if got := recorder.Body.String(); got != "messages are required\n" {
		t.Errorf(
			"body = %q, want %q",
			got,
			"messages are required\n",
		)
	}
}
