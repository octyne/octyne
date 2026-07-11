package openai

import "encoding/json"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StopSequences []string

func (s *StopSequences) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = StopSequences{single}
		return nil
	}

	var multiple []string
	if err := json.Unmarshal(data, &multiple); err != nil {
		return err
	}
	*s = multiple
	return nil
}

type LogitBias map[string]float64

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
	Metadata             *Metadata             `json:"metadata,omitempty"`
	ServiceTier          *ServiceTier          `json:"service_tier,omitempty"`
	PromptCacheOptions   *PromptCacheOptions   `json:"prompt_cache_options,omitempty"`
	Stop                 *StopSequences        `json:"stop,omitempty"`
	LogitBias            *LogitBias            `json:"logit_bias,omitempty"`
	StreamOptions        *StreamOptions        `json:"stream_options,omitempty"`
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

type Metadata map[string]string

type ServiceTier string

const (
	ServiceTierAuto     ServiceTier = "auto"
	ServiceTierDefault  ServiceTier = "default"
	ServiceTierFlex     ServiceTier = "flex"
	ServiceTierScale    ServiceTier = "scale"
	ServiceTierPriority ServiceTier = "priority"
)

type PromptCacheMode string

const (
	PromptCacheModeImplicit PromptCacheMode = "implicit"
	PromptCacheModeExplicit PromptCacheMode = "explicit"
)

type PromptCacheTTL string

const PromptCacheTTL30m PromptCacheTTL = "30m"

type PromptCacheOptions struct {
	Mode *PromptCacheMode `json:"mode,omitempty"`
	TTL  *PromptCacheTTL  `json:"ttl,omitempty"`
}

type StreamOptions struct {
	IncludeUsage       *bool `json:"include_usage,omitempty"`
	IncludeObfuscation *bool `json:"include_obfuscation,omitempty"`
}
