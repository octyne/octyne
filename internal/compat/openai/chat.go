package openai

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
}
