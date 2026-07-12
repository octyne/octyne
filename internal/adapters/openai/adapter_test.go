package openai

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/types"
)

func testTextChatMessage(role, text string) types.ChatMessage {
	return types.ChatMessage{
		Role:    role,
		Content: &types.MessageContent{Text: &text},
	}
}

func TestStreamChatDoesNotUseTotalRequestTimeout(t *testing.T) {
	const timeout = 100 * time.Millisecond

	upstream := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Error("upstream response does not support flushing")
				return
			}

			flusher.Flush()

			// This exceeds the non-streaming timeout after headers were received.
			time.Sleep(250 * time.Millisecond)

			_, err := io.WriteString(w, `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":123,"model":"gpt-5-nano","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null,"logprobs":null}]}

data: [DONE]

`)
			if err != nil {
				t.Errorf("write stream: %v", err)
			}
		},
	))
	defer upstream.Close()

	config := providers.Config{
		Name:                           "openai",
		BaseURL:                        upstream.URL,
		NonStreamingTimeout:            timeout,
		StreamingResponseHeaderTimeout: timeout,
	}

	adapter := New(config)

	chunks, err := adapter.StreamChat(
		context.Background(),
		types.ChatCompletionRequest{
			Model:    "gpt-5-nano",
			Messages: []types.ChatMessage{testTextChatMessage("user", "Hello")},
		},
	)
	if err != nil {
		t.Fatalf("StreamChat returned an error: %v", err)
	}

	chunk, ok := <-chunks
	if !ok {
		t.Fatal("stream closed before producing a chunk")
	}

	if chunk.Error != nil {
		t.Fatalf("unexpected stream error: %v", chunk.Error)
	}

	if len(chunk.Choices) != 1 ||
		chunk.Choices[0].Delta.Content == nil ||
		*chunk.Choices[0].Delta.Content != "Hello" {
		t.Fatalf("unexpected chunk: %+v", chunk)
	}

	if _, ok := <-chunks; ok {
		t.Fatal("stream remained open after [DONE]")
	}
}

func TestStreamChatReturnsUpstreamSetupError(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-request-id", "req_upstream")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = io.WriteString(w, `{"error":{"message":"rate limit exceeded","type":"rate_limit_error","param":null,"code":"rate_limit_exceeded"}}`)
		},
	))
	defer upstream.Close()

	adapter := New(providers.Config{
		Name:                           "openai",
		BaseURL:                        upstream.URL,
		NonStreamingTimeout:            time.Second,
		StreamingResponseHeaderTimeout: time.Second,
	})

	chunks, err := adapter.StreamChat(
		context.Background(),
		types.ChatCompletionRequest{
			Model:    "gpt-5-nano",
			Messages: []types.ChatMessage{testTextChatMessage("user", "Hello")},
		},
	)

	if err == nil {
		t.Fatal("StreamChat error = nil, want upstream status error")
	}

	if chunks != nil {
		t.Fatal("StreamChat returned chunks for failed setup")
	}

	var apiErr *types.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %T, want *types.APIError", err)
	}
	if apiErr.HTTPStatus != http.StatusTooManyRequests ||
		apiErr.Kind != types.ErrorKindRateLimit {
		t.Errorf("API error = %+v, want rate-limit status 429", apiErr)
	}
	if apiErr.ProviderRequestID != "req_upstream" {
		t.Errorf("ProviderRequestID = %q, want req_upstream", apiErr.ProviderRequestID)
	}
	if apiErr.Code == nil || *apiErr.Code != "rate_limit_exceeded" {
		t.Errorf("Code = %v, want rate_limit_exceeded", apiErr.Code)
	}
}

func TestChatSanitizesUpstreamServerError(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, `{"error":{"message":"secret internal provider detail","type":"server_error","param":null,"code":"internal"}}`)
	}))
	defer upstream.Close()

	adapter := New(providers.Config{
		Name:                           "openai",
		BaseURL:                        upstream.URL,
		NonStreamingTimeout:            time.Second,
		StreamingResponseHeaderTimeout: time.Second,
	})

	_, err := adapter.Chat(context.Background(), types.ChatCompletionRequest{
		Model:    "gpt-5-nano",
		Messages: []types.ChatMessage{testTextChatMessage("user", "Hello")},
	})
	var apiErr *types.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %T, want *types.APIError", err)
	}
	if apiErr.Message != "The upstream provider failed to process the request." {
		t.Errorf("Message = %q, want sanitized provider error", apiErr.Message)
	}
}

func TestUpstreamHTTPStatusMapping(t *testing.T) {
	tests := []struct {
		status int
		kind   types.ErrorKind
	}{
		{http.StatusBadRequest, types.ErrorKindInvalidRequest},
		{http.StatusUnauthorized, types.ErrorKindAuthentication},
		{http.StatusForbidden, types.ErrorKindPermission},
		{http.StatusNotFound, types.ErrorKindNotFound},
		{http.StatusConflict, types.ErrorKindInvalidRequest},
		{http.StatusUnprocessableEntity, types.ErrorKindInvalidRequest},
		{http.StatusTooManyRequests, types.ErrorKindRateLimit},
		{http.StatusGatewayTimeout, types.ErrorKindTimeout},
		{http.StatusServiceUnavailable, types.ErrorKindUnavailable},
		{http.StatusInternalServerError, types.ErrorKindInternal},
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			got := toAPIError(&http.Response{
				StatusCode: tt.status,
				Header:     make(http.Header),
			}, nil, nil)
			if got.Kind != tt.kind || got.HTTPStatus != tt.status {
				t.Errorf("error = %+v, want kind %q status %d", got, tt.kind, tt.status)
			}
		})
	}
}

func TestStreamChatPropagatesContextCancellation(t *testing.T) {
	requestCanceled := make(chan struct{})

	upstream := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Error("upstream response does not support flushing")
				return
			}

			flusher.Flush()
			<-r.Context().Done()
			close(requestCanceled)
		},
	))
	defer upstream.Close()

	adapter := New(providers.Config{
		Name:                           "openai",
		BaseURL:                        upstream.URL,
		NonStreamingTimeout:            time.Second,
		StreamingResponseHeaderTimeout: time.Second,
	})

	ctx, cancel := context.WithCancel(context.Background())

	chunks, err := adapter.StreamChat(
		ctx,
		types.ChatCompletionRequest{
			Model:    "gpt-5-nano",
			Messages: []types.ChatMessage{testTextChatMessage("user", "Hello")},
		},
	)
	if err != nil {
		cancel()
		t.Fatalf("StreamChat returned an error: %v", err)
	}

	cancel()

	select {
	case _, ok := <-chunks:
		if ok {
			t.Fatal("stream produced a chunk after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("stream channel did not close after cancellation")
	}

	select {
	case <-requestCanceled:
	case <-time.After(time.Second):
		t.Fatal("upstream request context was not cancelled")
	}
}
