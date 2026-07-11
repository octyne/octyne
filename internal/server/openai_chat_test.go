package server

import (
	"testing"

	openaicompat "github.com/octyne/octyne/internal/compat/openai"
)

func TestToCanonicalChatRequest(t *testing.T) {
	tests := []struct {
		name   string
		stream bool
	}{
		{name: "non-streaming", stream: false},
		{name: "streaming", stream: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := openaicompat.ChatCompletionRequest{
				Model: "gpt-5-nano",
				Messages: []openaicompat.Message{
					{
						Role:    "user",
						Content: "Hello",
					},
				},
				Stream: tt.stream,
			}

			got := toCanonicalChatRequest(req)

			if got.Model != req.Model {
				t.Errorf("Model = %q, want %q", got.Model, req.Model)
			}

			if got.Stream != req.Stream {
				t.Errorf("Stream = %t, want %t", got.Stream, req.Stream)
			}

			if len(got.Messages) != 1 {
				t.Fatalf(
					"len(Messages) = %d, want 1",
					len(got.Messages),
				)
			}

			if got.Messages[0].Role != "user" {
				t.Errorf(
					"Message.Role = %q, want user",
					got.Messages[0].Role,
				)
			}

			if got.Messages[0].Content != "Hello" {
				t.Errorf(
					"Message.Content = %q, want Hello",
					got.Messages[0].Content,
				)
			}
		})
	}
}
