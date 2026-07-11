package server

import (
	openaicompat "github.com/octyne/octyne/internal/compat/openai"
	"github.com/octyne/octyne/internal/types"
)

func toCanonicalChatRequest(
	req openaicompat.ChatCompletionRequest,
) types.ChatCompletionRequest {
	messages := make([]types.Message, len(req.Messages))

	for i, message := range req.Messages {
		messages[i] = types.Message{
			Role:    message.Role,
			Content: message.Content,
		}
	}
	return types.ChatCompletionRequest{
		Model:            req.Model,
		Messages:         messages,
		Stream:           req.Stream,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		FrequencyPenalty: req.FrequencyPenalty,
	}
}
