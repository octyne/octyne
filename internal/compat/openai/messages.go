package openai

import (
	"encoding/json"
	"fmt"
)

type MessageContent struct {
	Text  *string
	Parts *[]ContentPart
}

func (c *MessageContent) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		c.Text = &text
		c.Parts = nil
		return nil
	}
	var parts []ContentPart
	if err := json.Unmarshal(data, &parts); err != nil {
		return err
	}
	c.Text = nil
	c.Parts = &parts
	return nil
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

func (m *AssistantMessage) UnmarshalJSON(data []byte) error {
	type alias AssistantMessage
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*m = AssistantMessage(decoded)
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if content, ok := fields["content"]; ok && string(content) == "null" {
		m.ContentNull = true
	}
	return nil
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
type Message struct {
	Developer *DeveloperMessage
	System    *SystemMessage
	User      *UserMessage
	Assistant *AssistantMessage
	Tool      *ToolMessage
	Function  *FunctionMessage
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var h struct {
		Role string `json:"role"`
	}
	if err := json.Unmarshal(data, &h); err != nil {
		return err
	}
	switch h.Role {
	case "developer":
		m.Developer = &DeveloperMessage{}
		return json.Unmarshal(data, m.Developer)
	case "system":
		m.System = &SystemMessage{}
		return json.Unmarshal(data, m.System)
	case "user":
		m.User = &UserMessage{}
		return json.Unmarshal(data, m.User)
	case "assistant":
		m.Assistant = &AssistantMessage{}
		return json.Unmarshal(data, m.Assistant)
	case "tool":
		m.Tool = &ToolMessage{}
		return json.Unmarshal(data, m.Tool)
	case "function":
		m.Function = &FunctionMessage{}
		return json.Unmarshal(data, m.Function)
	default:
		return fmt.Errorf("unsupported chat message role %q", h.Role)
	}
}
