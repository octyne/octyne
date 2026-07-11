package server

import (
	"testing"

	openaicompat "github.com/octyne/octyne/internal/compat/openai"
)

func TestToCanonicalChatRequest(t *testing.T) {
	zeroTemperature := 0.0
	zeroTopP := 0.0
	zeroFrequencyPenalty := 0.0
	tests := []struct {
		name             string
		stream           bool
		temperature      *float64
		topP             *float64
		frequencyPenalty *float64
	}{
		{name: "omitted sampling fields", stream: false},
		{name: "explicit zero temperature", stream: true, temperature: &zeroTemperature},
		{name: "explicit zero top p", stream: false, topP: &zeroTopP},
		{
			name:             "explicit zero frequency penalty",
			stream:           true,
			frequencyPenalty: &zeroFrequencyPenalty,
		},
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
				Stream:           tt.stream,
				Temperature:      tt.temperature,
				TopP:             tt.topP,
				FrequencyPenalty: tt.frequencyPenalty,
			}

			got := toCanonicalChatRequest(req)

			if got.Model != req.Model {
				t.Errorf("Model = %q, want %q", got.Model, req.Model)
			}

			if got.Stream != req.Stream {
				t.Errorf("Stream = %t, want %t", got.Stream, req.Stream)
			}

			if tt.temperature == nil {
				if got.Temperature != nil {
					t.Errorf(
						"Temperature = %v, want nil",
						got.Temperature,
					)
				}
			} else {
				if got.Temperature == nil {
					t.Fatal("Temperature = nil, want explicit value")
				}

				if *got.Temperature != *tt.temperature {
					t.Errorf(
						"Temperature = %v, want %v",
						*got.Temperature,
						*tt.temperature,
					)
				}
			}

			if tt.topP == nil {
				if got.TopP != nil {
					t.Errorf(
						"TopP = %v, want nil",
						got.TopP,
					)
				}
			} else {
				if got.TopP == nil {
					t.Fatal("TopP = nil, want explicit value")
				}

				if *got.TopP != *tt.topP {
					t.Errorf(
						"TopP = %v, want %v",
						*got.TopP,
						*tt.topP,
					)
				}
			}

			if tt.frequencyPenalty == nil {
				if got.FrequencyPenalty != nil {
					t.Errorf(
						"FrequencyPenalty = %v, want nil",
						got.FrequencyPenalty,
					)
				}
			} else {
				if got.FrequencyPenalty == nil {
					t.Fatal("FrequencyPenalty = nil, want explicit value")
				}

				if *got.FrequencyPenalty != *tt.frequencyPenalty {
					t.Errorf(
						"FrequencyPenalty = %v, want %v",
						*got.FrequencyPenalty,
						*tt.frequencyPenalty,
					)
				}
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
