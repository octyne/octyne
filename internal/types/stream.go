package types

import "encoding/json"

type StreamFunctionCall struct {
	Arguments *string `json:"arguments,omitempty"`
	Name      *string `json:"name,omitempty"`
}

type StreamToolCall struct {
	Index    int                 `json:"index"`
	ID       *string             `json:"id,omitempty"`
	Function *StreamFunctionCall `json:"function,omitempty"`
	Type     *string             `json:"type,omitempty"`
}

type StreamDelta struct {
	Content      *string             `json:"content,omitempty"`
	FunctionCall *StreamFunctionCall `json:"function_call,omitempty"`
	Refusal      *string             `json:"refusal,omitempty"`
	Role         *string             `json:"role,omitempty"`
	ToolCalls    *[]StreamToolCall   `json:"tool_calls,omitempty"`
}

type StreamChoice struct {
	Index        int           `json:"index"`
	Delta        StreamDelta   `json:"delta"`
	FinishReason *FinishReason `json:"finish_reason"`
	Logprobs     *ChatLogprobs `json:"logprobs"`
}

type StreamChunk struct {
	ID                string                    `json:"id"`
	Object            string                    `json:"object"`
	Created           int64                     `json:"created"`
	Model             string                    `json:"model"`
	Choices           []StreamChoice            `json:"choices"`
	Moderation        *ChatCompletionModeration `json:"moderation,omitempty"`
	Obfuscation       *string                   `json:"obfuscation,omitempty"`
	ServiceTier       *ServiceTier              `json:"service_tier,omitempty"`
	SystemFingerprint *string                   `json:"system_fingerprint,omitempty"`
	Usage             *CompletionUsage          `json:"-"`
	UsagePresent      bool                      `json:"-"`
	Error             error                     `json:"-"`
}

func (c StreamChunk) MarshalJSON() ([]byte, error) {
	type wire struct {
		ID                string                    `json:"id"`
		Object            string                    `json:"object"`
		Created           int64                     `json:"created"`
		Model             string                    `json:"model"`
		Choices           []StreamChoice            `json:"choices"`
		Moderation        *ChatCompletionModeration `json:"moderation,omitempty"`
		Obfuscation       *string                   `json:"obfuscation,omitempty"`
		ServiceTier       *ServiceTier              `json:"service_tier,omitempty"`
		SystemFingerprint *string                   `json:"system_fingerprint,omitempty"`
		Usage             *CompletionUsage          `json:"usage"`
	}
	value := wire{
		ID: c.ID, Object: c.Object, Created: c.Created, Model: c.Model, Choices: c.Choices,
		Moderation: c.Moderation, Obfuscation: c.Obfuscation, ServiceTier: c.ServiceTier,
		SystemFingerprint: c.SystemFingerprint, Usage: c.Usage,
	}
	if c.UsagePresent {
		return json.Marshal(value)
	}
	type wireWithoutUsage struct {
		ID                string                    `json:"id"`
		Object            string                    `json:"object"`
		Created           int64                     `json:"created"`
		Model             string                    `json:"model"`
		Choices           []StreamChoice            `json:"choices"`
		Moderation        *ChatCompletionModeration `json:"moderation,omitempty"`
		Obfuscation       *string                   `json:"obfuscation,omitempty"`
		ServiceTier       *ServiceTier              `json:"service_tier,omitempty"`
		SystemFingerprint *string                   `json:"system_fingerprint,omitempty"`
	}
	return json.Marshal(wireWithoutUsage{
		ID: c.ID, Object: c.Object, Created: c.Created, Model: c.Model, Choices: c.Choices,
		Moderation: c.Moderation, Obfuscation: c.Obfuscation, ServiceTier: c.ServiceTier,
		SystemFingerprint: c.SystemFingerprint,
	})
}
