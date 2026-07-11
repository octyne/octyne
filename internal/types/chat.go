package types

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	Stream           bool      `json:"stream,omitempty"`
	Temperature      *float64  `json:"temperature,omitempty"`
	TopP             *float64  `json:"top_p,omitempty"`
	FrequencyPenalty *float64  `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64  `json:"presence_penalty,omitempty"`
	MaxOutputTokens  *int      `json:"max_output_tokens,omitempty"`
	CandidateCount   *int      `json:"candidate_count,omitempty"`
	ReturnLogprobs   *bool     `json:"return_logprobs,omitempty"`
	TopLogprobs      *int      `json:"top_logprobs,omitempty"`
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
