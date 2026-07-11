package openai

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model               string           `json:"model"`
	Messages            []Message        `json:"messages"`
	Stream              bool             `json:"stream,omitempty"`
	Temperature         *float64         `json:"temperature,omitempty"`
	TopP                *float64         `json:"top_p,omitempty"`
	FrequencyPenalty    *float64         `json:"frequency_penalty,omitempty"`
	PresencePenalty     *float64         `json:"presence_penalty,omitempty"`
	MaxCompletionTokens *int             `json:"max_completion_tokens,omitempty"`
	N                   *int             `json:"n,omitempty"`
	Logprobs            *bool            `json:"logprobs,omitempty"`
	TopLogprobs         *int             `json:"top_logprobs,omitempty"`
	ReasoningEffort     *ReasoningEffort `json:"reasoning_effort,omitempty"`
	Verbosity           *Verbosity       `json:"verbosity,omitempty"`
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
