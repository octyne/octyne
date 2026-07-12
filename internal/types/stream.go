package types

type StreamDelta struct {
	Role    *string `json:"role,omitempty"`
	Content *string `json:"content,omitempty"`
}

type StreamChoice struct {
	Index        int           `json:"index"`
	Delta        StreamDelta   `json:"delta"`
	FinishReason *FinishReason `json:"finish_reason"`
	Logprobs     *ChatLogprobs `json:"logprobs"`
}

type StreamChunk struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []StreamChoice   `json:"choices"`
	Usage   *CompletionUsage `json:"usage,omitempty"`
	Error   error            `json:"-"`
}
