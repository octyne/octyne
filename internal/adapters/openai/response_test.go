package openai

import (
	"encoding/json"
	"testing"
)

func TestChatCompletionResponseTranslatesTypedAccounting(t *testing.T) {
	input := []byte(`{
		"id":"chatcmpl-123",
		"object":"chat.completion",
		"created":123,
		"model":"gpt-5-nano",
		"choices":[{
			"index":0,
			"message":{"role":"assistant","content":"Hello"},
			"finish_reason":"length",
			"logprobs":{
				"content":[{
					"token":"Hello",
					"bytes":[72,101,108,108,111],
					"logprob":-0.25,
					"top_logprobs":[{"token":"Hi","bytes":null,"logprob":-1.5}]
				}],
				"refusal":null
			}
		}],
		"usage":{
			"completion_tokens":7,
			"prompt_tokens":5,
			"total_tokens":12,
			"completion_tokens_details":{
				"accepted_prediction_tokens":0,
				"audio_tokens":2,
				"reasoning_tokens":3,
				"rejected_prediction_tokens":1
			},
			"prompt_tokens_details":{
				"audio_tokens":0,
				"cache_write_tokens":4,
				"cached_tokens":2
			}
		}
	}`)

	var providerResponse ChatCompletionResponse
	if err := json.Unmarshal(input, &providerResponse); err != nil {
		t.Fatalf("decode provider response: %v", err)
	}

	response := toChatCompletionResponse(providerResponse)
	choice := response.Choices[0]
	if choice.FinishReason == nil || *choice.FinishReason != "length" {
		t.Fatalf("FinishReason = %v, want length", choice.FinishReason)
	}
	if choice.Logprobs == nil || len(choice.Logprobs.Content) != 1 {
		t.Fatalf("Logprobs = %+v, want one content token", choice.Logprobs)
	}
	token := choice.Logprobs.Content[0]
	if token.Token != "Hello" || len(token.Bytes) != 5 || len(token.TopLogprobs) != 1 ||
		token.TopLogprobs[0].Bytes != nil {
		t.Errorf("unexpected token log probabilities: %+v", token)
	}
	if choice.Logprobs.Refusal != nil {
		t.Errorf("Refusal logprobs = %v, want null", choice.Logprobs.Refusal)
	}

	if response.Usage == nil || response.Usage.CompletionTokens != 7 ||
		response.Usage.PromptTokens != 5 || response.Usage.TotalTokens != 12 {
		t.Fatalf("Usage = %+v, want typed token totals", response.Usage)
	}
	completionDetails := response.Usage.CompletionTokensDetails
	if completionDetails == nil || completionDetails.AcceptedPredictionTokens == nil ||
		*completionDetails.AcceptedPredictionTokens != 0 || completionDetails.AudioTokens == nil ||
		*completionDetails.AudioTokens != 2 || completionDetails.ReasoningTokens == nil ||
		*completionDetails.ReasoningTokens != 3 || completionDetails.RejectedPredictionTokens == nil ||
		*completionDetails.RejectedPredictionTokens != 1 {
		t.Errorf("CompletionTokensDetails = %+v", completionDetails)
	}
	promptDetails := response.Usage.PromptTokensDetails
	if promptDetails == nil || promptDetails.AudioTokens == nil || *promptDetails.AudioTokens != 0 ||
		promptDetails.CacheWriteTokens == nil || *promptDetails.CacheWriteTokens != 4 ||
		promptDetails.CachedTokens == nil || *promptDetails.CachedTokens != 2 {
		t.Errorf("PromptTokensDetails = %+v", promptDetails)
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("encode canonical response: %v", err)
	}
	var got any
	var want any
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatalf("decode encoded response: %v", err)
	}
	if err := json.Unmarshal(input, &want); err != nil {
		t.Fatalf("decode expected response: %v", err)
	}
	if !jsonEqual(got, want) {
		t.Errorf("translated response:\n%s\nwant:\n%s", encoded, input)
	}
}

func jsonEqual(left, right any) bool {
	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)
	return string(leftJSON) == string(rightJSON)
}
