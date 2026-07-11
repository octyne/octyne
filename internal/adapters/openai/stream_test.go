package openai

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestReadChatCompletionStream(t *testing.T) {
	body := io.NopCloser(strings.NewReader(
		`data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":123,"model":"gpt-5-nano","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null,"logprobs":null}]}

data: [DONE]

`,
	))

	chunks := readChatCompletionStream(
		context.Background(),
		body,
	)

	chunk, ok := <-chunks
	if !ok {
		t.Fatal("stream closed before producing a chunk")
	}

	if chunk.Error != nil {
		t.Fatalf("unexpected stream error: %v", chunk.Error)
	}

	if chunk.ID != "chatcmpl-123" {
		t.Errorf("ID = %q, want %q", chunk.ID, "chatcmpl-123")
	}

	if chunk.Object != "chat.completion.chunk" {
		t.Errorf(
			"Object = %q, want %q",
			chunk.Object,
			"chat.completion.chunk",
		)
	}

	if len(chunk.Choices) != 1 {
		t.Fatalf(
			"len(Choices) = %d, want 1",
			len(chunk.Choices),
		)
	}

	choice := chunk.Choices[0]

	if choice.Delta.Role == nil ||
		*choice.Delta.Role != "assistant" {
		t.Errorf("Delta.Role = %v, want assistant", choice.Delta.Role)
	}

	if choice.Delta.Content == nil ||
		*choice.Delta.Content != "Hello" {
		t.Errorf("Delta.Content = %v, want Hello", choice.Delta.Content)
	}

	if _, ok := <-chunks; ok {
		t.Fatal("stream remained open after [DONE]")
	}
}

func TestReadChatCompletionStreamReportsInvalidJSON(
	t *testing.T,
) {
	body := io.NopCloser(strings.NewReader(
		"data: {invalid-json}\n\n",
	))

	chunks := readChatCompletionStream(
		context.Background(),
		body,
	)

	chunk, ok := <-chunks
	if !ok {
		t.Fatal("stream closed without reporting the error")
	}

	if chunk.Error == nil {
		t.Fatal("Error = nil, want JSON decoding error")
	}

	if !strings.Contains(
		chunk.Error.Error(),
		"decode openai stream chunk",
	) {
		t.Errorf(
			"Error = %q, want decoding context",
			chunk.Error,
		)
	}

	if _, ok := <-chunks; ok {
		t.Fatal("stream remained open after decoding error")
	}
}

func TestReadChatCompletionStreamClosesBodyExactlyOnce(t *testing.T) {
	body := &countingReadCloser{
		Reader: strings.NewReader("data: [DONE]\n\n"),
	}

	chunks := readChatCompletionStream(
		context.Background(),
		body,
	)

	for chunk := range chunks {
		t.Fatalf("unexpected chunk: %+v", chunk)
	}

	if body.closeCount != 1 {
		t.Errorf(
			"body close count = %d, want 1",
			body.closeCount,
		)
	}
}

type countingReadCloser struct {
	io.Reader
	closeCount int
}

func (r *countingReadCloser) Close() error {
	r.closeCount++
	return nil
}
