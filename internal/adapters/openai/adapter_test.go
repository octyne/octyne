package openai

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/types"
)

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
			Model: "gpt-5-nano",
			Messages: []types.Message{
				{
					Role:    "user",
					Content: "Hello",
				},
			},
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
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
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
			Model: "gpt-5-nano",
			Messages: []types.Message{
				{Role: "user", Content: "Hello"},
			},
		},
	)

	if err == nil {
		t.Fatal("StreamChat error = nil, want upstream status error")
	}

	if chunks != nil {
		t.Fatal("StreamChat returned chunks for failed setup")
	}

	if !strings.Contains(err.Error(), "status 429") {
		t.Errorf("error = %q, want status 429", err)
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
			Model: "gpt-5-nano",
			Messages: []types.Message{
				{Role: "user", Content: "Hello"},
			},
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
