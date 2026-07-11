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
		Metadata:               toCanonicalMetadata(req.Metadata),
		ServiceTier:            toCanonicalServiceTier(req.ServiceTier),
		PromptCacheOptions:     toCanonicalPromptCacheOptions(req.PromptCacheOptions),
		StopSequences:          toCanonicalStopSequences(req.Stop),
		LogitBias:              toCanonicalLogitBias(req.LogitBias),
		StreamOptions:          toCanonicalStreamOptions(req.StreamOptions),
		Modalities:             toCanonicalModalities(req.Modalities),
		AudioOutput:            toCanonicalAudioOutput(req.Audio),
		ResponseFormat:         toCanonicalResponseFormat(req.ResponseFormat),
		Prediction:             toCanonicalPrediction(req.Prediction),
		Moderation:             toCanonicalModeration(req.Moderation),
		WebSearch:              toCanonicalWebSearch(req.WebSearchOptions),
	}
}

func toCanonicalModeration(value *openaicompat.ModerationOptions) *types.ModerationOptions {
	if value == nil {
		return nil
	}
	converted := &types.ModerationOptions{Model: value.Model}
	if value.Policy != nil {
		converted.Policy = &types.ModerationPolicy{}
		if value.Policy.Input != nil {
			converted.Policy.Input = &types.ModerationRule{Mode: types.ModerationMode(value.Policy.Input.Mode)}
		}
		if value.Policy.Output != nil {
			converted.Policy.Output = &types.ModerationRule{Mode: types.ModerationMode(value.Policy.Output.Mode)}
		}
	}
	return converted
}

func toCanonicalWebSearch(value *openaicompat.WebSearchOptions) *types.WebSearchOptions {
	if value == nil {
		return nil
	}
	converted := &types.WebSearchOptions{}
	if value.SearchContextSize != nil {
		size := types.SearchContextSize(*value.SearchContextSize)
		converted.SearchContextSize = &size
	}
	if value.UserLocation != nil {
		converted.UserLocation = &types.UserLocation{
			Type: value.UserLocation.Type,
			Approximate: types.ApproximateLocation{
				City: value.UserLocation.Approximate.City, Country: value.UserLocation.Approximate.Country,
				Region: value.UserLocation.Approximate.Region, Timezone: value.UserLocation.Approximate.Timezone,
			},
		}
	}
	return converted
}

func toCanonicalPrediction(value *openaicompat.Prediction) *types.Prediction {
	if value == nil {
		return nil
	}
	converted := &types.Prediction{Type: value.Type}
	if value.Content.Text != nil {
		converted.Content.Text = value.Content.Text
	}
	if value.Content.Parts != nil {
		parts := make([]types.TextContentPart, len(*value.Content.Parts))
		for i, part := range *value.Content.Parts {
			parts[i] = types.TextContentPart{Type: part.Type, Text: part.Text}
			if part.PromptCacheBreakpoint != nil {
				parts[i].PromptCacheBreakpoint = &types.PromptCacheBreakpoint{
					Mode: part.PromptCacheBreakpoint.Mode,
				}
			}
		}
		converted.Content.Parts = &parts
	}
	return converted
}

func toCanonicalResponseFormat(value *openaicompat.ResponseFormat) *types.ResponseFormat {
	if value == nil {
		return nil
	}
	converted := &types.ResponseFormat{Type: types.ResponseFormatType(value.Type)}
	if value.JSONSchema != nil {
		converted.JSONSchema = &types.JSONSchemaFormat{
			Name:        value.JSONSchema.Name,
			Description: value.JSONSchema.Description,
			Schema:      value.JSONSchema.Schema,
			Strict:      value.JSONSchema.Strict,
		}
	}
	return converted
}

func toCanonicalModalities(value *openaicompat.Modalities) *types.Modalities {
	if value == nil {
		return nil
	}
	converted := make(types.Modalities, len(*value))
	for i, modality := range *value {
		converted[i] = types.Modality(modality)
	}
	return &converted
}

func toCanonicalAudioOutput(value *openaicompat.AudioOutput) *types.AudioOutput {
	if value == nil {
		return nil
	}
	return &types.AudioOutput{
		Format: types.AudioFormat(value.Format),
		Voice: types.AudioVoice{
			Name: value.Voice.Name,
			ID:   value.Voice.ID,
		},
	}
}

func toCanonicalStreamOptions(value *openaicompat.StreamOptions) *types.StreamOptions {
	if value == nil {
		return nil
	}
	return &types.StreamOptions{
		IncludeUsage:       value.IncludeUsage,
		IncludeObfuscation: value.IncludeObfuscation,
	}
}

func toCanonicalStopSequences(value *openaicompat.StopSequences) *types.StopSequences {
	if value == nil {
		return nil
	}
	converted := types.StopSequences(*value)
	return &converted
}

func toCanonicalLogitBias(value *openaicompat.LogitBias) *types.LogitBias {
	if value == nil {
		return nil
	}
	converted := types.LogitBias(*value)
	return &converted
}

func toCanonicalMetadata(value *openaicompat.Metadata) *types.Metadata {
	if value == nil {
		return nil
	}
	converted := types.Metadata(*value)
	return &converted
}

func toCanonicalServiceTier(value *openaicompat.ServiceTier) *types.ServiceTier {
	if value == nil {
		return nil
	}
	converted := types.ServiceTier(*value)
	return &converted
}

func toCanonicalPromptCacheOptions(
	value *openaicompat.PromptCacheOptions,
) *types.PromptCacheOptions {
	if value == nil {
		return nil
	}
	converted := &types.PromptCacheOptions{}
	if value.Mode != nil {
		mode := types.PromptCacheMode(*value.Mode)
		converted.Mode = &mode
	}
	if value.TTL != nil {
		ttl := types.PromptCacheTTL(*value.TTL)
		converted.TTL = &ttl
	}
	return converted
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
