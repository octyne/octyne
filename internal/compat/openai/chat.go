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
	Name *string
	ID   *string
}

func (v *AudioVoice) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err == nil {
		v.Name = &name
		v.ID = nil
		return nil
	}
	var custom struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &custom); err != nil {
		return err
	}
	v.Name = nil
	v.ID = &custom.ID
	return nil
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
	Text  *string
	Parts *[]TextContentPart
}

func (c *PredictionContent) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		c.Text = &text
		c.Parts = nil
		return nil
	}
	var parts []TextContentPart
	if err := json.Unmarshal(data, &parts); err != nil {
		return err
	}
	c.Text = nil
	c.Parts = &parts
	return nil
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
	Model                string                      `json:"model"`
	Messages             []Message                   `json:"messages"`
	Stream               bool                        `json:"stream,omitempty"`
	Temperature          *float64                    `json:"temperature,omitempty"`
	TopP                 *float64                    `json:"top_p,omitempty"`
	FrequencyPenalty     *float64                    `json:"frequency_penalty,omitempty"`
	PresencePenalty      *float64                    `json:"presence_penalty,omitempty"`
	MaxCompletionTokens  *int                        `json:"max_completion_tokens,omitempty"`
	N                    *int                        `json:"n,omitempty"`
	Logprobs             *bool                       `json:"logprobs,omitempty"`
	TopLogprobs          *int                        `json:"top_logprobs,omitempty"`
	ReasoningEffort      *ReasoningEffort            `json:"reasoning_effort,omitempty"`
	Verbosity            *Verbosity                  `json:"verbosity,omitempty"`
	Seed                 *int64                      `json:"seed,omitempty"`
	Store                *bool                       `json:"store,omitempty"`
	ParallelToolCalls    *bool                       `json:"parallel_tool_calls,omitempty"`
	SafetyIdentifier     *string                     `json:"safety_identifier,omitempty"`
	PromptCacheKey       *string                     `json:"prompt_cache_key,omitempty"`
	MaxTokens            *int                        `json:"max_tokens,omitempty"`
	User                 *string                     `json:"user,omitempty"`
	PromptCacheRetention *PromptCacheRetention       `json:"prompt_cache_retention,omitempty"`
	Metadata             *Metadata                   `json:"metadata,omitempty"`
	ServiceTier          *ServiceTier                `json:"service_tier,omitempty"`
	PromptCacheOptions   *PromptCacheOptions         `json:"prompt_cache_options,omitempty"`
	Stop                 *StopSequences              `json:"stop,omitempty"`
	LogitBias            *LogitBias                  `json:"logit_bias,omitempty"`
	StreamOptions        *StreamOptions              `json:"stream_options,omitempty"`
	Modalities           *Modalities                 `json:"modalities,omitempty"`
	Audio                *AudioOutput                `json:"audio,omitempty"`
	ResponseFormat       *ResponseFormat             `json:"response_format,omitempty"`
	Prediction           *Prediction                 `json:"prediction,omitempty"`
	Moderation           *ModerationOptions          `json:"moderation,omitempty"`
	WebSearchOptions     *WebSearchOptions           `json:"web_search_options,omitempty"`
	Tools                *[]Tool                     `json:"tools,omitempty"`
	ToolChoice           *ToolChoice                 `json:"tool_choice,omitempty"`
	Functions            *[]LegacyFunctionDefinition `json:"functions,omitempty"`
	FunctionCall         *LegacyFunctionCall         `json:"function_call,omitempty"`
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
