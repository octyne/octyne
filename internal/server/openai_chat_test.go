package server

import (
	"testing"

	openaicompat "github.com/octyne/octyne/internal/compat/openai"
)

func TestToCanonicalChatRequest(t *testing.T) {
	zeroTemperature := 0.0
	zeroTopP := 0.0
	zeroFrequencyPenalty := 0.0
	zeroPresencePenalty := 0.0
	zeroMaxCompletionTokens := 0
	zeroN := 0
	falseLogprobs := false
	zeroTopLogprobs := 0
	tests := []struct {
		name                string
		stream              bool
		temperature         *float64
		topP                *float64
		frequencyPenalty    *float64
		presencePenalty     *float64
		maxCompletionTokens *int
		n                   *int
		logprobs            *bool
		topLogprobs         *int
	}{
		{name: "omitted sampling fields", stream: false},
		{name: "explicit zero temperature", stream: true, temperature: &zeroTemperature},
		{name: "explicit zero top p", stream: false, topP: &zeroTopP},
		{
			name:             "explicit zero frequency penalty",
			stream:           true,
			frequencyPenalty: &zeroFrequencyPenalty,
		},
		{
			name:            "explicit zero presence penalty",
			presencePenalty: &zeroPresencePenalty,
		},
		{
			name:                "explicit zero max completion tokens",
			maxCompletionTokens: &zeroMaxCompletionTokens,
		},
		{name: "explicit zero n", n: &zeroN},
		{name: "explicit false logprobs", logprobs: &falseLogprobs},
		{
			name:        "explicit zero top logprobs",
			topLogprobs: &zeroTopLogprobs,
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
				Stream:              tt.stream,
				Temperature:         tt.temperature,
				TopP:                tt.topP,
				FrequencyPenalty:    tt.frequencyPenalty,
				PresencePenalty:     tt.presencePenalty,
				MaxCompletionTokens: tt.maxCompletionTokens,
				N:                   tt.n,
				Logprobs:            tt.logprobs,
				TopLogprobs:         tt.topLogprobs,
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

			if tt.presencePenalty == nil {
				if got.PresencePenalty != nil {
					t.Errorf(
						"PresencePenalty = %v, want nil",
						got.PresencePenalty,
					)
				}
			} else if got.PresencePenalty == nil {
				t.Fatal("PresencePenalty = nil, want explicit value")
			} else if *got.PresencePenalty != *tt.presencePenalty {
				t.Errorf(
					"PresencePenalty = %v, want %v",
					*got.PresencePenalty,
					*tt.presencePenalty,
				)
			}

			if tt.maxCompletionTokens == nil {
				if got.MaxOutputTokens != nil {
					t.Errorf(
						"MaxOutputTokens = %v, want nil",
						got.MaxOutputTokens,
					)
				}
			} else if got.MaxOutputTokens == nil {
				t.Fatal("MaxOutputTokens = nil, want explicit value")
			} else if *got.MaxOutputTokens != *tt.maxCompletionTokens {
				t.Errorf(
					"MaxOutputTokens = %v, want %v",
					*got.MaxOutputTokens,
					*tt.maxCompletionTokens,
				)
			}

			if tt.n == nil {
				if got.CandidateCount != nil {
					t.Errorf(
						"CandidateCount = %v, want nil",
						got.CandidateCount,
					)
				}
			} else if got.CandidateCount == nil {
				t.Fatal("CandidateCount = nil, want explicit value")
			} else if *got.CandidateCount != *tt.n {
				t.Errorf(
					"CandidateCount = %v, want %v",
					*got.CandidateCount,
					*tt.n,
				)
			}

			if tt.logprobs == nil {
				if got.ReturnLogprobs != nil {
					t.Errorf(
						"ReturnLogprobs = %v, want nil",
						got.ReturnLogprobs,
					)
				}
			} else if got.ReturnLogprobs == nil {
				t.Fatal("ReturnLogprobs = nil, want explicit value")
			} else if *got.ReturnLogprobs != *tt.logprobs {
				t.Errorf(
					"ReturnLogprobs = %t, want %t",
					*got.ReturnLogprobs,
					*tt.logprobs,
				)
			}

			if tt.topLogprobs == nil {
				if got.TopLogprobs != nil {
					t.Errorf(
						"TopLogprobs = %v, want nil",
						got.TopLogprobs,
					)
				}
			} else if got.TopLogprobs == nil {
				t.Fatal("TopLogprobs = nil, want explicit value")
			} else if *got.TopLogprobs != *tt.topLogprobs {
				t.Errorf(
					"TopLogprobs = %d, want %d",
					*got.TopLogprobs,
					*tt.topLogprobs,
				)
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
