package server

import (
	openaicompat "github.com/octyne/octyne/internal/compat/openai"
	"github.com/octyne/octyne/internal/types"
)

func toCanonicalMessages(messages []openaicompat.Message) []types.ChatMessage {
	converted := make([]types.ChatMessage, len(messages))
	for i, message := range messages {
		switch {
		case message.Developer != nil:
			converted[i] = types.ChatMessage{Role: "developer", Content: toCanonicalMessageContent(&message.Developer.Content), Name: message.Developer.Name}
		case message.System != nil:
			converted[i] = types.ChatMessage{Role: "system", Content: toCanonicalMessageContent(&message.System.Content), Name: message.System.Name}
		case message.User != nil:
			converted[i] = types.ChatMessage{Role: "user", Content: toCanonicalMessageContent(&message.User.Content), Name: message.User.Name}
		case message.Assistant != nil:
			m := message.Assistant
			converted[i] = types.ChatMessage{Role: "assistant", Content: toCanonicalMessageContent(m.Content), Name: m.Name, Refusal: m.Refusal, ToolCalls: toCanonicalMessageToolCalls(m.ToolCalls), FunctionCall: toCanonicalMessageFunctionCall(m.FunctionCall)}
			converted[i].ContentNull = m.ContentNull
			if m.Audio != nil {
				converted[i].Audio = &types.AudioReference{ID: m.Audio.ID}
			}
		case message.Tool != nil:
			id := message.Tool.ToolCallID
			converted[i] = types.ChatMessage{Role: "tool", Content: toCanonicalMessageContent(&message.Tool.Content), ToolCallID: &id}
		case message.Function != nil:
			content := message.Function.Content
			name := message.Function.Name
			converted[i] = types.ChatMessage{Role: "function", Content: &types.MessageContent{Text: &content}, Name: &name}
		}
	}
	return converted
}

func toCanonicalMessageContent(value *openaicompat.MessageContent) *types.MessageContent {
	if value == nil {
		return nil
	}
	content := &types.MessageContent{Text: value.Text}
	if value.Parts != nil {
		parts := make([]types.ContentPart, len(*value.Parts))
		for i, part := range *value.Parts {
			parts[i] = types.ContentPart{Type: part.Type, Text: part.Text, Refusal: part.Refusal}
			if part.ImageURL != nil {
				parts[i].ImageURL = &types.ImageURL{URL: part.ImageURL.URL, Detail: part.ImageURL.Detail}
			}
			if part.InputAudio != nil {
				parts[i].InputAudio = &types.InputAudio{Data: part.InputAudio.Data, Format: part.InputAudio.Format}
			}
			if part.File != nil {
				parts[i].File = &types.FileInput{FileData: part.File.FileData, FileID: part.File.FileID, Filename: part.File.Filename}
			}
			if part.PromptCacheBreakpoint != nil {
				parts[i].PromptCacheBreakpoint = &types.PromptCacheBreakpoint{Mode: part.PromptCacheBreakpoint.Mode}
			}
		}
		content.Parts = &parts
	}
	return content
}

func toCanonicalMessageFunctionCall(value *openaicompat.MessageFunctionCall) *types.MessageFunctionCall {
	if value == nil {
		return nil
	}
	return &types.MessageFunctionCall{Arguments: value.Arguments, Name: value.Name}
}

func toCanonicalMessageToolCalls(value *[]openaicompat.MessageToolCall) *[]types.MessageToolCall {
	if value == nil {
		return nil
	}
	calls := make([]types.MessageToolCall, len(*value))
	for i, call := range *value {
		calls[i] = types.MessageToolCall{ID: call.ID, Type: call.Type, Function: toCanonicalMessageFunctionCall(call.Function)}
		if call.Custom != nil {
			calls[i].Custom = &types.MessageCustomCall{Input: call.Custom.Input, Name: call.Custom.Name}
		}
	}
	return &calls
}
