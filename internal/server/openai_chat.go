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
		Model:                  req.Model,
		Messages:               messages,
		Stream:                 req.Stream,
		Temperature:            req.Temperature,
		TopP:                   req.TopP,
		FrequencyPenalty:       req.FrequencyPenalty,
		PresencePenalty:        req.PresencePenalty,
		MaxOutputTokens:        req.MaxCompletionTokens,
		CandidateCount:         req.N,
		ReturnLogprobs:         req.Logprobs,
		TopLogprobs:            req.TopLogprobs,
		ReasoningEffort:        toCanonicalReasoningEffort(req.ReasoningEffort),
		Verbosity:              toCanonicalVerbosity(req.Verbosity),
		Seed:                   req.Seed,
		StoreOutput:            req.Store,
		AllowParallelToolCalls: req.ParallelToolCalls,
		SafetyIdentifier:       req.SafetyIdentifier,
		PromptCacheKey:         req.PromptCacheKey,
		LegacyMaxOutputTokens:  req.MaxTokens,
		LegacyUser:             req.User,
		PromptCacheRetention:   toCanonicalPromptCacheRetention(req.PromptCacheRetention),
	}
}

func toCanonicalPromptCacheRetention(
	value *openaicompat.PromptCacheRetention,
) *types.PromptCacheRetention {
	if value == nil {
		return nil
	}
	converted := types.PromptCacheRetention(*value)
	return &converted
}

func toCanonicalReasoningEffort(
	value *openaicompat.ReasoningEffort,
) *types.ReasoningEffort {
	if value == nil {
		return nil
	}

	converted := types.ReasoningEffort(*value)
	return &converted
}

func toCanonicalVerbosity(
	value *openaicompat.Verbosity,
) *types.Verbosity {
	if value == nil {
		return nil
	}

	converted := types.Verbosity(*value)
	return &converted
}
