package openai

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ReasoningEffort string

const (
	ReasoningEffortNone    ReasoningEffort = "none"
	ReasoningEffortMinimal ReasoningEffort = "minimal"
	ReasoningEffortLow     ReasoningEffort = "low"
	ReasoningEffortMedium  ReasoningEffort = "medium"
	ReasoningEffortHigh    ReasoningEffort = "high"
	ReasoningEffortXHigh   ReasoningEffort = "xhigh"
	ReasoningEffortMax     ReasoningEffort = "max"
)

type Verbosity string

const (
	VerbosityLow    Verbosity = "low"
	VerbosityMedium Verbosity = "medium"
	VerbosityHigh   Verbosity = "high"
)

type PromptCacheRetention string

const (
	PromptCacheRetentionInMemory PromptCacheRetention = "in_memory"
	PromptCacheRetention24h      PromptCacheRetention = "24h"
)

type ChatCompletionRequest struct {
	Model                string                `json:"model"`
	Messages             []Message             `json:"messages"`
	Stream               bool                  `json:"stream,omitempty"`
	Temperature          *float64              `json:"temperature,omitempty"`
	TopP                 *float64              `json:"top_p,omitempty"`
	FrequencyPenalty     *float64              `json:"frequency_penalty,omitempty"`
	PresencePenalty      *float64              `json:"presence_penalty,omitempty"`
	MaxCompletionTokens  *int                  `json:"max_completion_tokens,omitempty"`
	N                    *int                  `json:"n,omitempty"`
	Logprobs             *bool                 `json:"logprobs,omitempty"`
	TopLogprobs          *int                  `json:"top_logprobs,omitempty"`
	ReasoningEffort      *ReasoningEffort      `json:"reasoning_effort,omitempty"`
	Verbosity            *Verbosity            `json:"verbosity,omitempty"`
	Seed                 *int64                `json:"seed,omitempty"`
	Store                *bool                 `json:"store,omitempty"`
	ParallelToolCalls    *bool                 `json:"parallel_tool_calls,omitempty"`
	SafetyIdentifier     *string               `json:"safety_identifier,omitempty"`
	PromptCacheKey       *string               `json:"prompt_cache_key,omitempty"`
	MaxTokens            *int                  `json:"max_tokens,omitempty"`
	User                 *string               `json:"user,omitempty"`
	PromptCacheRetention *PromptCacheRetention `json:"prompt_cache_retention,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason *string `json:"finish_reason"`
	Logprobs     any     `json:"logprobs"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   any      `json:"usage,omitempty"`
}

type Delta struct {
	Role    *string `json:"role,omitempty"`
	Content *string `json:"content,omitempty"`
}

type ChunkChoice struct {
	Index        int     `json:"index"`
	Delta        Delta   `json:"delta"`
	FinishReason *string `json:"finish_reason"`
	Logprobs     any     `json:"logprobs"`
}

type ChatCompletionChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChunkChoice `json:"choices"`
	Usage   any           `json:"usage,omitempty"`
}
