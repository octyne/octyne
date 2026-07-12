package types

type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"
	FinishReasonLength        FinishReason = "length"
	FinishReasonToolCalls     FinishReason = "tool_calls"
	FinishReasonContentFilter FinishReason = "content_filter"
	FinishReasonFunctionCall  FinishReason = "function_call"
)

type URLCitation struct {
	EndIndex   int    `json:"end_index"`
	StartIndex int    `json:"start_index"`
	Title      string `json:"title"`
	URL        string `json:"url"`
}

type Annotation struct {
	Type        string      `json:"type"`
	URLCitation URLCitation `json:"url_citation"`
}

type ChatCompletionAudio struct {
	ID         string `json:"id"`
	Data       string `json:"data"`
	ExpiresAt  int64  `json:"expires_at"`
	Transcript string `json:"transcript"`
}

type ResponseFunctionCall struct {
	Arguments string `json:"arguments"`
	Name      string `json:"name"`
}

type ResponseCustomCall struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

type ResponseToolCall struct {
	ID       string                `json:"id"`
	Type     string                `json:"type"`
	Function *ResponseFunctionCall `json:"function,omitempty"`
	Custom   *ResponseCustomCall   `json:"custom,omitempty"`
}

type ResponseMessage struct {
	Role         string                `json:"role"`
	Content      *string               `json:"content"`
	Refusal      *string               `json:"refusal"`
	Annotations  *[]Annotation         `json:"annotations,omitempty"`
	Audio        *ChatCompletionAudio  `json:"audio,omitempty"`
	FunctionCall *ResponseFunctionCall `json:"function_call,omitempty"`
	ToolCalls    *[]ResponseToolCall   `json:"tool_calls,omitempty"`
}

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
