package types

import "encoding/json"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StopSequences []string

type LogitBias map[string]float64

type Modality string

const (
	ModalityText  Modality = "text"
	ModalityAudio Modality = "audio"
)

type Modalities []Modality

type AudioFormat string

const (
	AudioFormatWAV   AudioFormat = "wav"
	AudioFormatAAC   AudioFormat = "aac"
	AudioFormatMP3   AudioFormat = "mp3"
	AudioFormatFLAC  AudioFormat = "flac"
	AudioFormatOpus  AudioFormat = "opus"
	AudioFormatPCM16 AudioFormat = "pcm16"
)

type AudioVoice struct {
	Name *string `json:"name,omitempty"`
	ID   *string `json:"id,omitempty"`
}

type AudioOutput struct {
	Format AudioFormat `json:"format"`
	Voice  AudioVoice  `json:"voice"`
}

type ResponseFormatType string

const (
	ResponseFormatText       ResponseFormatType = "text"
	ResponseFormatJSONObject ResponseFormatType = "json_object"
	ResponseFormatJSONSchema ResponseFormatType = "json_schema"
)

type JSONSchemaFormat struct {
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	Schema      *json.RawMessage `json:"schema,omitempty"`
	Strict      *bool            `json:"strict,omitempty"`
}

type ResponseFormat struct {
	Type       ResponseFormatType `json:"type"`
	JSONSchema *JSONSchemaFormat  `json:"json_schema,omitempty"`
}

type PromptCacheBreakpoint struct {
	Mode string `json:"mode"`
}

type TextContentPart struct {
	Type                  string                 `json:"type"`
	Text                  string                 `json:"text"`
	PromptCacheBreakpoint *PromptCacheBreakpoint `json:"prompt_cache_breakpoint,omitempty"`
}

type PredictionContent struct {
	Text  *string            `json:"text,omitempty"`
	Parts *[]TextContentPart `json:"parts,omitempty"`
}

type Prediction struct {
	Type    string            `json:"type"`
	Content PredictionContent `json:"content"`
}

type ModerationMode string

const (
	ModerationModeScore ModerationMode = "score"
	ModerationModeBlock ModerationMode = "block"
)

type ModerationRule struct {
	Mode ModerationMode `json:"mode"`
}

type ModerationPolicy struct {
	Input  *ModerationRule `json:"input,omitempty"`
	Output *ModerationRule `json:"output,omitempty"`
}

type ModerationOptions struct {
	Model  string            `json:"model"`
	Policy *ModerationPolicy `json:"policy,omitempty"`
}

type SearchContextSize string

const (
	SearchContextSizeLow    SearchContextSize = "low"
	SearchContextSizeMedium SearchContextSize = "medium"
	SearchContextSizeHigh   SearchContextSize = "high"
)

type ApproximateLocation struct {
	City     *string `json:"city,omitempty"`
	Country  *string `json:"country,omitempty"`
	Region   *string `json:"region,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
}

type UserLocation struct {
	Type        string              `json:"type"`
	Approximate ApproximateLocation `json:"approximate"`
}

type WebSearchOptions struct {
	SearchContextSize *SearchContextSize `json:"search_context_size,omitempty"`
	UserLocation      *UserLocation      `json:"user_location,omitempty"`
}

type ChatCompletionRequest struct {
	Model                  string                      `json:"model"`
	Messages               []ChatMessage               `json:"messages"`
	Stream                 bool                        `json:"stream,omitempty"`
	Temperature            *float64                    `json:"temperature,omitempty"`
	TopP                   *float64                    `json:"top_p,omitempty"`
	FrequencyPenalty       *float64                    `json:"frequency_penalty,omitempty"`
	PresencePenalty        *float64                    `json:"presence_penalty,omitempty"`
	MaxOutputTokens        *int                        `json:"max_output_tokens,omitempty"`
	CandidateCount         *int                        `json:"candidate_count,omitempty"`
	ReturnLogprobs         *bool                       `json:"return_logprobs,omitempty"`
	TopLogprobs            *int                        `json:"top_logprobs,omitempty"`
	ReasoningEffort        *ReasoningEffort            `json:"reasoning_effort,omitempty"`
	Verbosity              *Verbosity                  `json:"verbosity,omitempty"`
	Seed                   *int64                      `json:"seed,omitempty"`
	StoreOutput            *bool                       `json:"store_output,omitempty"`
	AllowParallelToolCalls *bool                       `json:"allow_parallel_tool_calls,omitempty"`
	SafetyIdentifier       *string                     `json:"safety_identifier,omitempty"`
	PromptCacheKey         *string                     `json:"prompt_cache_key,omitempty"`
	LegacyMaxOutputTokens  *int                        `json:"legacy_max_output_tokens,omitempty"`
	LegacyUser             *string                     `json:"legacy_user,omitempty"`
	PromptCacheRetention   *PromptCacheRetention       `json:"prompt_cache_retention,omitempty"`
	Metadata               *Metadata                   `json:"metadata,omitempty"`
	ServiceTier            *ServiceTier                `json:"service_tier,omitempty"`
	PromptCacheOptions     *PromptCacheOptions         `json:"prompt_cache_options,omitempty"`
	StopSequences          *StopSequences              `json:"stop_sequences,omitempty"`
	LogitBias              *LogitBias                  `json:"logit_bias,omitempty"`
	StreamOptions          *StreamOptions              `json:"stream_options,omitempty"`
	Modalities             *Modalities                 `json:"modalities,omitempty"`
	AudioOutput            *AudioOutput                `json:"audio_output,omitempty"`
	ResponseFormat         *ResponseFormat             `json:"response_format,omitempty"`
	Prediction             *Prediction                 `json:"prediction,omitempty"`
	Moderation             *ModerationOptions          `json:"moderation,omitempty"`
	WebSearch              *WebSearchOptions           `json:"web_search,omitempty"`
	Tools                  *[]Tool                     `json:"tools,omitempty"`
	ToolChoice             *ToolChoice                 `json:"tool_choice,omitempty"`
	LegacyFunctions        *[]LegacyFunctionDefinition `json:"legacy_functions,omitempty"`
	LegacyFunctionCall     *LegacyFunctionCall         `json:"legacy_function_call,omitempty"`
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

type StreamOptions struct {
	IncludeUsage       *bool `json:"include_usage,omitempty"`
	IncludeObfuscation *bool `json:"include_obfuscation,omitempty"`
}
