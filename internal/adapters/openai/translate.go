package openai

import "github.com/octyne/octyne/internal/types"

func toChatCompletionRequest(
	req types.ChatCompletionRequest,
	stream bool,
) ChatCompletionRequest {
	messages := make([]Message, 0, len(req.Messages))

	for _, msg := range req.Messages {
		messages = append(messages, Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Stream:      stream,
		Temperature: req.Temperature,
	}
}

func toChatCompletionResponse(
	resp ChatCompletionResponse,
) types.ChatCompletionResponse {
	choices := make([]types.Choice, 0, len(resp.Choices))

	for _, choice := range resp.Choices {
		choices = append(choices, types.Choice{
			Index: choice.Index,
			Message: types.Message{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
			FinishReason: choice.FinishReason,
			Logprobs:     choice.Logprobs,
		})
	}

	return types.ChatCompletionResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: choices,
		Usage:   resp.Usage,
	}
}

func toStreamChunk(
	chunk ChatCompletionChunk,
) types.StreamChunk {
	choices := make([]types.StreamChoice, 0, len(chunk.Choices))

	for _, choice := range chunk.Choices {
		choices = append(choices, types.StreamChoice{
			Index: choice.Index,
			Delta: types.StreamDelta{
				Role:    choice.Delta.Role,
				Content: choice.Delta.Content,
			},
			FinishReason: choice.FinishReason,
			Logprobs:     choice.Logprobs,
		})
	}

	return types.StreamChunk{
		ID:      chunk.ID,
		Object:  chunk.Object,
		Created: chunk.Created,
		Model:   chunk.Model,
		Choices: choices,
		Usage:   chunk.Usage,
	}
}
