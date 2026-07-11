package types

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StopSequences []string

type LogitBias map[string]float64

type ChatCompletionRequest struct {
	Model                  string                `json:"model"`
	Messages               []Message             `json:"messages"`
	Stream                 bool                  `json:"stream,omitempty"`
	Temperature            *float64              `json:"temperature,omitempty"`
	TopP                   *float64              `json:"top_p,omitempty"`
	FrequencyPenalty       *float64              `json:"frequency_penalty,omitempty"`
	PresencePenalty        *float64              `json:"presence_penalty,omitempty"`
	MaxOutputTokens        *int                  `json:"max_output_tokens,omitempty"`
	CandidateCount         *int                  `json:"candidate_count,omitempty"`
	ReturnLogprobs         *bool                 `json:"return_logprobs,omitempty"`
	TopLogprobs            *int                  `json:"top_logprobs,omitempty"`
	ReasoningEffort        *ReasoningEffort      `json:"reasoning_effort,omitempty"`
	Verbosity              *Verbosity            `json:"verbosity,omitempty"`
	Seed                   *int64                `json:"seed,omitempty"`
	StoreOutput            *bool                 `json:"store_output,omitempty"`
	AllowParallelToolCalls *bool                 `json:"allow_parallel_tool_calls,omitempty"`
	SafetyIdentifier       *string               `json:"safety_identifier,omitempty"`
	PromptCacheKey         *string               `json:"prompt_cache_key,omitempty"`
	LegacyMaxOutputTokens  *int                  `json:"legacy_max_output_tokens,omitempty"`
	LegacyUser             *string               `json:"legacy_user,omitempty"`
	PromptCacheRetention   *PromptCacheRetention `json:"prompt_cache_retention,omitempty"`
	Metadata               *Metadata             `json:"metadata,omitempty"`
	ServiceTier            *ServiceTier          `json:"service_tier,omitempty"`
	PromptCacheOptions     *PromptCacheOptions   `json:"prompt_cache_options,omitempty"`
	StopSequences          *StopSequences        `json:"stop_sequences,omitempty"`
	LogitBias              *LogitBias            `json:"logit_bias,omitempty"`
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
