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
		Model:                req.Model,
		Messages:             messages,
		Stream:               stream,
		Temperature:          req.Temperature,
		TopP:                 req.TopP,
		FrequencyPenalty:     req.FrequencyPenalty,
		PresencePenalty:      req.PresencePenalty,
		MaxCompletionTokens:  req.MaxOutputTokens,
		N:                    req.CandidateCount,
		Logprobs:             req.ReturnLogprobs,
		TopLogprobs:          req.TopLogprobs,
		ReasoningEffort:      toReasoningEffort(req.ReasoningEffort),
		Verbosity:            toVerbosity(req.Verbosity),
		Seed:                 req.Seed,
		Store:                req.StoreOutput,
		ParallelToolCalls:    req.AllowParallelToolCalls,
		SafetyIdentifier:     req.SafetyIdentifier,
		PromptCacheKey:       req.PromptCacheKey,
		MaxTokens:            req.LegacyMaxOutputTokens,
		User:                 req.LegacyUser,
		PromptCacheRetention: toPromptCacheRetention(req.PromptCacheRetention),
		Metadata:             toMetadata(req.Metadata),
		ServiceTier:          toServiceTier(req.ServiceTier),
		PromptCacheOptions:   toPromptCacheOptions(req.PromptCacheOptions),
		Stop:                 toStopSequences(req.StopSequences),
		LogitBias:            toLogitBias(req.LogitBias),
		StreamOptions:        toStreamOptions(req.StreamOptions),
	}
}

func toStreamOptions(value *types.StreamOptions) *StreamOptions {
	if value == nil {
		return nil
	}
	return &StreamOptions{
		IncludeUsage:       value.IncludeUsage,
		IncludeObfuscation: value.IncludeObfuscation,
	}
}

func toStopSequences(value *types.StopSequences) *StopSequences {
	if value == nil {
		return nil
	}
	converted := StopSequences(*value)
	return &converted
}

func toLogitBias(value *types.LogitBias) *LogitBias {
	if value == nil {
		return nil
	}
	converted := LogitBias(*value)
	return &converted
}

func toMetadata(value *types.Metadata) *Metadata {
	if value == nil {
		return nil
	}
	converted := Metadata(*value)
	return &converted
}

func toServiceTier(value *types.ServiceTier) *ServiceTier {
	if value == nil {
		return nil
	}
	converted := ServiceTier(*value)
	return &converted
}

func toPromptCacheOptions(value *types.PromptCacheOptions) *PromptCacheOptions {
	if value == nil {
		return nil
	}
	converted := &PromptCacheOptions{}
	if value.Mode != nil {
		mode := PromptCacheMode(*value.Mode)
		converted.Mode = &mode
	}
	if value.TTL != nil {
		ttl := PromptCacheTTL(*value.TTL)
		converted.TTL = &ttl
	}
	return converted
}

func toPromptCacheRetention(
	value *types.PromptCacheRetention,
) *PromptCacheRetention {
	if value == nil {
		return nil
	}
	converted := PromptCacheRetention(*value)
	return &converted
}

func toReasoningEffort(value *types.ReasoningEffort) *ReasoningEffort {
	if value == nil {
		return nil
	}

	converted := ReasoningEffort(*value)
	return &converted
}

func toVerbosity(value *types.Verbosity) *Verbosity {
	if value == nil {
		return nil
	}

	converted := Verbosity(*value)
	return &converted
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
