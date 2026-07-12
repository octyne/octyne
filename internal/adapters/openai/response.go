package openai

import (
	"encoding/json"
	"fmt"
)

type FinishReason string

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

type ResponseMessage struct {
	Role         string               `json:"role"`
	Content      *string              `json:"content"`
	Refusal      *string              `json:"refusal"`
	Annotations  *[]Annotation        `json:"annotations,omitempty"`
	Audio        *ChatCompletionAudio `json:"audio,omitempty"`
	FunctionCall *MessageFunctionCall `json:"function_call,omitempty"`
	ToolCalls    *[]MessageToolCall   `json:"tool_calls,omitempty"`
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

type ModerationResult struct {
	Categories                map[string]bool     `json:"categories"`
	CategoryAppliedInputTypes map[string][]string `json:"category_applied_input_types"`
	CategoryScores            map[string]float64  `json:"category_scores"`
	Flagged                   bool                `json:"flagged"`
	Model                     string              `json:"model"`
	Type                      string              `json:"type"`
}

type ModerationResults struct {
	Model   string             `json:"model"`
	Results []ModerationResult `json:"results"`
	Type    string             `json:"type"`
}

type ModerationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type ModerationOutcome struct {
	Results *ModerationResults
	Error   *ModerationError
}

func (o *ModerationOutcome) UnmarshalJSON(data []byte) error {
	var discriminator struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &discriminator); err != nil {
		return err
	}
	switch discriminator.Type {
	case "moderation_results":
		var value ModerationResults
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		o.Results = &value
		o.Error = nil
		return nil
	case "error":
		var value ModerationError
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		o.Error = &value
		o.Results = nil
		return nil
	default:
		return fmt.Errorf("unsupported moderation outcome type %q", discriminator.Type)
	}
}

type ChatCompletionModeration struct {
	Input  ModerationOutcome `json:"input"`
	Output ModerationOutcome `json:"output"`
}

type StreamFunctionCall struct {
	Arguments *string `json:"arguments,omitempty"`
	Name      *string `json:"name,omitempty"`
}

type StreamToolCall struct {
	Index    int                 `json:"index"`
	ID       *string             `json:"id,omitempty"`
	Function *StreamFunctionCall `json:"function,omitempty"`
	Type     *string             `json:"type,omitempty"`
}

func (c *ChatCompletionChunk) UnmarshalJSON(data []byte) error {
	type alias ChatCompletionChunk
	wire := struct {
		Usage json.RawMessage `json:"usage"`
		*alias
	}{alias: (*alias)(c)}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	c.UsagePresent = len(wire.Usage) != 0
	c.Usage = nil
	if !c.UsagePresent || string(wire.Usage) == "null" {
		return nil
	}
	return json.Unmarshal(wire.Usage, &c.Usage)
}
