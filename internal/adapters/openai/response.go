package openai

type FinishReason string

type TopLogprob struct {
	Token   string  `json:"token"`
	Bytes   []int   `json:"bytes"`
	Logprob float64 `json:"logprob"`
}

type TokenLogprob struct {
	Token       string       `json:"token"`
	Bytes       []int        `json:"bytes"`
	Logprob     float64      `json:"logprob"`
	TopLogprobs []TopLogprob `json:"top_logprobs"`
}

type ChatLogprobs struct {
	Content []TokenLogprob `json:"content"`
	Refusal []TokenLogprob `json:"refusal"`
}

type CompletionTokensDetails struct {
	AcceptedPredictionTokens *int `json:"accepted_prediction_tokens,omitempty"`
	AudioTokens              *int `json:"audio_tokens,omitempty"`
	ReasoningTokens          *int `json:"reasoning_tokens,omitempty"`
	RejectedPredictionTokens *int `json:"rejected_prediction_tokens,omitempty"`
}

type PromptTokensDetails struct {
	AudioTokens      *int `json:"audio_tokens,omitempty"`
	CacheWriteTokens *int `json:"cache_write_tokens,omitempty"`
	CachedTokens     *int `json:"cached_tokens,omitempty"`
}

type CompletionUsage struct {
	CompletionTokens        int                      `json:"completion_tokens"`
	PromptTokens            int                      `json:"prompt_tokens"`
	TotalTokens             int                      `json:"total_tokens"`
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`
}
