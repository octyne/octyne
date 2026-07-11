package openai

import "encoding/json"

type MessageContent struct {
	Text  *string
	Parts *[]ContentPart
}

func (c MessageContent) MarshalJSON() ([]byte, error) {
	if c.Text != nil {
		return json.Marshal(*c.Text)
	}
	if c.Parts != nil {
		return json.Marshal(*c.Parts)
	}
	return []byte("null"), nil
}

type ImageURL struct {
	URL    string  `json:"url"`
	Detail *string `json:"detail,omitempty"`
}
type InputAudio struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}
type FileInput struct {
	FileData *string `json:"file_data,omitempty"`
	FileID   *string `json:"file_id,omitempty"`
	Filename *string `json:"filename,omitempty"`
}
type ContentPart struct {
	Type                  string                 `json:"type"`
	Text                  *string                `json:"text,omitempty"`
	ImageURL              *ImageURL              `json:"image_url,omitempty"`
	InputAudio            *InputAudio            `json:"input_audio,omitempty"`
	File                  *FileInput             `json:"file,omitempty"`
	Refusal               *string                `json:"refusal,omitempty"`
	PromptCacheBreakpoint *PromptCacheBreakpoint `json:"prompt_cache_breakpoint,omitempty"`
}
type AudioReference struct {
	ID string `json:"id"`
}
type MessageFunctionCall struct {
	Arguments string `json:"arguments"`
	Name      string `json:"name"`
}
type MessageCustomCall struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}
type MessageToolCall struct {
	ID       string               `json:"id"`
	Type     string               `json:"type"`
	Function *MessageFunctionCall `json:"function,omitempty"`
	Custom   *MessageCustomCall   `json:"custom,omitempty"`
}
type DeveloperMessage struct {
	Content MessageContent `json:"content"`
	Role    string         `json:"role"`
	Name    *string        `json:"name,omitempty"`
}
type SystemMessage struct {
	Content MessageContent `json:"content"`
	Role    string         `json:"role"`
	Name    *string        `json:"name,omitempty"`
}
type UserMessage struct {
	Content MessageContent `json:"content"`
	Role    string         `json:"role"`
	Name    *string        `json:"name,omitempty"`
}
type AssistantMessage struct {
	Role         string               `json:"role"`
	Audio        *AudioReference      `json:"audio,omitempty"`
	Content      *MessageContent      `json:"content,omitempty"`
	FunctionCall *MessageFunctionCall `json:"function_call,omitempty"`
	Name         *string              `json:"name,omitempty"`
	Refusal      *string              `json:"refusal,omitempty"`
	ToolCalls    *[]MessageToolCall   `json:"tool_calls,omitempty"`
	ContentNull  bool                 `json:"-"`
}

func (m AssistantMessage) MarshalJSON() ([]byte, error) {
	type wire struct {
		Role         string               `json:"role"`
		Audio        *AudioReference      `json:"audio,omitempty"`
		Content      *json.RawMessage     `json:"content,omitempty"`
		FunctionCall *MessageFunctionCall `json:"function_call,omitempty"`
		Name         *string              `json:"name,omitempty"`
		Refusal      *string              `json:"refusal,omitempty"`
		ToolCalls    *[]MessageToolCall   `json:"tool_calls,omitempty"`
	}
	result := wire{Role: m.Role, Audio: m.Audio, FunctionCall: m.FunctionCall, Name: m.Name, Refusal: m.Refusal, ToolCalls: m.ToolCalls}
	if m.Content != nil {
		encoded, err := json.Marshal(m.Content)
		if err != nil {
			return nil, err
		}
		raw := json.RawMessage(encoded)
		result.Content = &raw
	} else if m.ContentNull {
		raw := json.RawMessage("null")
		result.Content = &raw
	}
	return json.Marshal(result)
}

type ToolMessage struct {
	Content    MessageContent `json:"content"`
	Role       string         `json:"role"`
	ToolCallID string         `json:"tool_call_id"`
}
type FunctionMessage struct {
	Content string `json:"content"`
	Name    string `json:"name"`
	Role    string `json:"role"`
}
type RequestMessage struct {
	Developer *DeveloperMessage
	System    *SystemMessage
	User      *UserMessage
	Assistant *AssistantMessage
	Tool      *ToolMessage
	Function  *FunctionMessage
}

func (m RequestMessage) MarshalJSON() ([]byte, error) {
	switch {
	case m.Developer != nil:
		return json.Marshal(m.Developer)
	case m.System != nil:
		return json.Marshal(m.System)
	case m.User != nil:
		return json.Marshal(m.User)
	case m.Assistant != nil:
		return json.Marshal(m.Assistant)
	case m.Tool != nil:
		return json.Marshal(m.Tool)
	case m.Function != nil:
		return json.Marshal(m.Function)
	default:
		return []byte("null"), nil
	}
}
