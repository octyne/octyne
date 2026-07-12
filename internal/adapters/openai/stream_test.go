package openai

import (
	"context"
	"encoding/json"
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

func TestReadChatCompletionStreamPreservesTypedDeltasAndMetadata(t *testing.T) {
	input := `data: {"id":"chatcmpl-stream","object":"chat.completion.chunk","created":123,"model":"gpt-5-nano","choices":[{"index":0,"delta":{"role":"assistant","refusal":"cannot","function_call":{"name":"legacy","arguments":"{"},"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"weather","arguments":"{\"city\":"}}]},"finish_reason":null,"logprobs":{"content":null,"refusal":[{"token":"cannot","bytes":null,"logprob":-0.1,"top_logprobs":[]}]}},{"index":1,"delta":{"content":"alternative"},"finish_reason":"stop","logprobs":null}],"moderation":{"input":{"code":"unavailable","message":"try again","type":"error"},"output":{"code":"unavailable","message":"try again","type":"error"}},"obfuscation":"random-padding","service_tier":"flex","system_fingerprint":"fp_123","usage":null}

data: {"id":"chatcmpl-stream","object":"chat.completion.chunk","created":123,"model":"gpt-5-nano","choices":[],"usage":{"completion_tokens":3,"prompt_tokens":2,"total_tokens":5}}

data: [DONE]

`
	chunks := readChatCompletionStream(
		context.Background(),
		io.NopCloser(strings.NewReader(input)),
	)

	first, ok := <-chunks
	if !ok {
		t.Fatal("stream closed before first chunk")
	}
	if first.Error != nil {
		t.Fatalf("first chunk error: %v", first.Error)
	}
	if !first.UsagePresent || first.Usage != nil {
		t.Errorf("first usage = %+v, present = %t; want explicit null", first.Usage, first.UsagePresent)
	}
	if first.ServiceTier == nil || *first.ServiceTier != "flex" || first.SystemFingerprint == nil ||
		*first.SystemFingerprint != "fp_123" || first.Obfuscation == nil || *first.Obfuscation != "random-padding" {
		t.Errorf("unexpected stream metadata: %+v", first)
	}
	if first.Moderation == nil || first.Moderation.Input.Error == nil {
		t.Errorf("Moderation = %+v, want typed errors", first.Moderation)
	}
	delta := first.Choices[0].Delta
	if delta.Refusal == nil || *delta.Refusal != "cannot" || delta.FunctionCall == nil ||
		delta.FunctionCall.Name == nil || *delta.FunctionCall.Name != "legacy" || delta.ToolCalls == nil ||
		len(*delta.ToolCalls) != 1 || (*delta.ToolCalls)[0].Function == nil {
		t.Errorf("Delta = %+v, want refusal and function fragments", delta)
	}
	if len(first.Choices) != 2 || first.Choices[0].Logprobs == nil ||
		len(first.Choices[0].Logprobs.Refusal) != 1 || first.Choices[1].FinishReason == nil ||
		*first.Choices[1].FinishReason != "stop" {
		t.Errorf("Choices = %+v, want multiple typed choices and refusal logprobs", first.Choices)
	}
	encodedFirst, err := json.Marshal(first)
	if err != nil {
		t.Fatalf("encode first chunk: %v", err)
	}
	if !strings.Contains(string(encodedFirst), `"usage":null`) {
		t.Errorf("first chunk does not preserve explicit null usage: %s", encodedFirst)
	}

	last, ok := <-chunks
	if !ok {
		t.Fatal("stream closed before usage chunk")
	}
	if last.Error != nil {
		t.Fatalf("usage chunk error: %v", last.Error)
	}
	if !last.UsagePresent || last.Usage == nil || last.Usage.TotalTokens != 5 || len(last.Choices) != 0 {
		t.Errorf("usage chunk = %+v, want empty choices and token totals", last)
	}
	encodedLast, err := json.Marshal(last)
	if err != nil {
		t.Fatalf("encode usage chunk: %v", err)
	}
	if !strings.Contains(string(encodedLast), `"choices":[],"usage":{"completion_tokens":3`) {
		t.Errorf("unexpected usage-only chunk: %s", encodedLast)
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
