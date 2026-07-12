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
			"message":{"role":"assistant","content":"Hello","refusal":null},
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

func TestChatCompletionResponseTranslatesAssistantOutputs(t *testing.T) {
	input := []byte(`{
		"id":"chatcmpl-outputs",
		"object":"chat.completion",
		"created":456,
		"model":"gpt-5-nano",
		"choices":[
			{
				"index":0,
				"message":{
					"role":"assistant",
					"content":null,
					"refusal":"cannot comply",
					"annotations":[{
						"type":"url_citation",
						"url_citation":{"end_index":12,"start_index":0,"title":"Example","url":"https://example.com"}
					}],
					"audio":{"id":"audio_1","data":"AA==","expires_at":789,"transcript":"spoken answer"},
					"tool_calls":[
						{"id":"call_1","type":"function","function":{"arguments":"{\"city\":\"Pune\"}","name":"weather"}},
						{"id":"call_2","type":"custom","custom":{"input":"pwd","name":"shell"}}
					],
					"function_call":{"arguments":"{}","name":"legacy"}
				},
				"finish_reason":"tool_calls",
				"logprobs":null
			}
		],
		"usage":{"completion_tokens":1,"prompt_tokens":2,"total_tokens":3}
	}`)

	var providerResponse ChatCompletionResponse
	if err := json.Unmarshal(input, &providerResponse); err != nil {
		t.Fatalf("decode provider response: %v", err)
	}
	response := toChatCompletionResponse(providerResponse)
	message := response.Choices[0].Message
	if message.Content != nil || message.Refusal == nil || *message.Refusal != "cannot comply" {
		t.Errorf("unexpected nullable output: %+v", message)
	}
	if message.Annotations == nil || len(*message.Annotations) != 1 ||
		(*message.Annotations)[0].URLCitation.URL != "https://example.com" {
		t.Errorf("Annotations = %+v", message.Annotations)
	}
	if message.Audio == nil || message.Audio.ID != "audio_1" || message.Audio.Transcript != "spoken answer" {
		t.Errorf("Audio = %+v", message.Audio)
	}
	if message.ToolCalls == nil || len(*message.ToolCalls) != 2 ||
		(*message.ToolCalls)[0].Function == nil || (*message.ToolCalls)[1].Custom == nil {
		t.Errorf("ToolCalls = %+v", message.ToolCalls)
	}
	if message.FunctionCall == nil || message.FunctionCall.Name != "legacy" {
		t.Errorf("FunctionCall = %+v", message.FunctionCall)
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

func TestChatCompletionResponseTranslatesTopLevelMetadata(t *testing.T) {
	input := []byte(`{
		"id":"chatcmpl-metadata",
		"object":"chat.completion",
		"created":789,
		"model":"gpt-5-nano",
		"choices":[],
		"moderation":{
			"input":{
				"model":"omni-moderation-latest",
				"results":[{
					"categories":{"violence":false},
					"category_applied_input_types":{"violence":["text"]},
					"category_scores":{"violence":0.01},
					"flagged":false,
					"model":"omni-moderation-latest",
					"type":"moderation_result"
				}],
				"type":"moderation_results"
			},
			"output":{"code":"moderation_unavailable","message":"try again","type":"error"}
		},
		"service_tier":"priority",
		"system_fingerprint":"fp_123",
		"usage":{"completion_tokens":0,"prompt_tokens":0,"total_tokens":0}
	}`)

	var providerResponse ChatCompletionResponse
	if err := json.Unmarshal(input, &providerResponse); err != nil {
		t.Fatalf("decode provider response: %v", err)
	}
	response := toChatCompletionResponse(providerResponse)
	if response.ServiceTier == nil || *response.ServiceTier != "priority" {
		t.Errorf("ServiceTier = %v, want priority", response.ServiceTier)
	}
	if response.SystemFingerprint == nil || *response.SystemFingerprint != "fp_123" {
		t.Errorf("SystemFingerprint = %v, want fp_123", response.SystemFingerprint)
	}
	if response.Moderation == nil || response.Moderation.Input.Results == nil ||
		len(response.Moderation.Input.Results.Results) != 1 || response.Moderation.Output.Error == nil ||
		response.Moderation.Output.Error.Code != "moderation_unavailable" {
		t.Errorf("Moderation = %+v", response.Moderation)
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
