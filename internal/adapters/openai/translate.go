package openai

import "github.com/usekeel/keel/internal/types"

func toChatCompletionRequest(
	req types.ChatCompletionRequest,
) ChatCompletionRequest {
	messages := make([]Message, 0, len(req.Messages))

	for _, msg := range req.Messages {
		messages = append(messages, Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return ChatCompletionRequest{
		Model:    req.Model,
		Messages: messages,
	}
}
