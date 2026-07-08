package openai

import "github.com/octyne/octyne/internal/types"

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

func toChatCompletionResponse(
	resp ChatCompletionResponse,
) types.ChatCompletionResponse {
	choices := make([]types.Choice, 0, len(resp.Choices))

	for _, choice := range resp.Choices {
		choices = append(choices, types.Choice{
			Message: types.Message{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
		})
	}

	return types.ChatCompletionResponse{
		ID:      resp.ID,
		Model:   resp.Model,
		Choices: choices,
	}
}
