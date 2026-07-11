package openai

import "github.com/octyne/octyne/internal/types"

func toRequestMessages(messages []types.ChatMessage) []RequestMessage {
	converted := make([]RequestMessage, len(messages))
	for i, message := range messages {
		switch message.Role {
		case "developer":
			converted[i].Developer = &DeveloperMessage{Role: "developer", Content: messageContentValue(message.Content), Name: message.Name}
		case "system":
			converted[i].System = &SystemMessage{Role: "system", Content: messageContentValue(message.Content), Name: message.Name}
		case "user":
			converted[i].User = &UserMessage{Role: "user", Content: messageContentValue(message.Content), Name: message.Name}
		case "assistant":
			converted[i].Assistant = &AssistantMessage{Role: "assistant", Content: toMessageContent(message.Content), Name: message.Name, Refusal: message.Refusal, ToolCalls: toMessageToolCalls(message.ToolCalls), FunctionCall: toMessageFunctionCall(message.FunctionCall)}
			converted[i].Assistant.ContentNull = message.ContentNull
			if message.Audio != nil {
				converted[i].Assistant.Audio = &AudioReference{ID: message.Audio.ID}
			}
		case "tool":
			id := ""
			if message.ToolCallID != nil {
				id = *message.ToolCallID
			}
			converted[i].Tool = &ToolMessage{Role: "tool", Content: messageContentValue(message.Content), ToolCallID: id}
		case "function":
			content, name := "", ""
			if message.Content != nil && message.Content.Text != nil {
				content = *message.Content.Text
			}
			if message.Name != nil {
				name = *message.Name
			}
			converted[i].Function = &FunctionMessage{Role: "function", Content: content, Name: name}
		}
	}
	return converted
}

func messageContentValue(value *types.MessageContent) MessageContent {
	converted := toMessageContent(value)
	if converted == nil {
		return MessageContent{}
	}
	return *converted
}
func toMessageContent(value *types.MessageContent) *MessageContent {
	if value == nil {
		return nil
	}
	content := &MessageContent{Text: value.Text}
	if value.Parts != nil {
		parts := make([]ContentPart, len(*value.Parts))
		for i, part := range *value.Parts {
			parts[i] = ContentPart{Type: part.Type, Text: part.Text, Refusal: part.Refusal}
			if part.ImageURL != nil {
				parts[i].ImageURL = &ImageURL{URL: part.ImageURL.URL, Detail: part.ImageURL.Detail}
			}
			if part.InputAudio != nil {
				parts[i].InputAudio = &InputAudio{Data: part.InputAudio.Data, Format: part.InputAudio.Format}
			}
			if part.File != nil {
				parts[i].File = &FileInput{FileData: part.File.FileData, FileID: part.File.FileID, Filename: part.File.Filename}
			}
			if part.PromptCacheBreakpoint != nil {
				parts[i].PromptCacheBreakpoint = &PromptCacheBreakpoint{Mode: part.PromptCacheBreakpoint.Mode}
			}
		}
		content.Parts = &parts
	}
	return content
}
func toMessageFunctionCall(value *types.MessageFunctionCall) *MessageFunctionCall {
	if value == nil {
		return nil
	}
	return &MessageFunctionCall{Arguments: value.Arguments, Name: value.Name}
}
func toMessageToolCalls(value *[]types.MessageToolCall) *[]MessageToolCall {
	if value == nil {
		return nil
	}
	calls := make([]MessageToolCall, len(*value))
	for i, call := range *value {
		calls[i] = MessageToolCall{ID: call.ID, Type: call.Type, Function: toMessageFunctionCall(call.Function)}
		if call.Custom != nil {
			calls[i].Custom = &MessageCustomCall{Input: call.Custom.Input, Name: call.Custom.Name}
		}
	}
	return &calls
}
