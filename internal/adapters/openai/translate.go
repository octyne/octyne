package openai

import "github.com/octyne/octyne/internal/types"

func toChatCompletionRequest(
	req types.ChatCompletionRequest,
	stream bool,
) ChatCompletionRequest {
	return ChatCompletionRequest{
		Model:                req.Model,
		Messages:             toRequestMessages(req.Messages),
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
		Modalities:           toModalities(req.Modalities),
		Audio:                toAudioOutput(req.AudioOutput),
		ResponseFormat:       toResponseFormat(req.ResponseFormat),
		Prediction:           toPrediction(req.Prediction),
		Moderation:           toModeration(req.Moderation),
		WebSearchOptions:     toWebSearchOptions(req.WebSearch),
		Tools:                toTools(req.Tools),
		ToolChoice:           toToolChoice(req.ToolChoice),
		Functions:            toLegacyFunctions(req.LegacyFunctions),
		FunctionCall:         toLegacyFunctionCall(req.LegacyFunctionCall),
	}
}

func toModeration(value *types.ModerationOptions) *ModerationOptions {
	if value == nil {
		return nil
	}
	converted := &ModerationOptions{Model: value.Model}
	if value.Policy != nil {
		converted.Policy = &ModerationPolicy{}
		if value.Policy.Input != nil {
			converted.Policy.Input = &ModerationRule{Mode: ModerationMode(value.Policy.Input.Mode)}
		}
		if value.Policy.Output != nil {
			converted.Policy.Output = &ModerationRule{Mode: ModerationMode(value.Policy.Output.Mode)}
		}
	}
	return converted
}

func toWebSearchOptions(value *types.WebSearchOptions) *WebSearchOptions {
	if value == nil {
		return nil
	}
	converted := &WebSearchOptions{}
	if value.SearchContextSize != nil {
		size := SearchContextSize(*value.SearchContextSize)
		converted.SearchContextSize = &size
	}
	if value.UserLocation != nil {
		converted.UserLocation = &UserLocation{
			Type: value.UserLocation.Type,
			Approximate: ApproximateLocation{
				City: value.UserLocation.Approximate.City, Country: value.UserLocation.Approximate.Country,
				Region: value.UserLocation.Approximate.Region, Timezone: value.UserLocation.Approximate.Timezone,
			},
		}
	}
	return converted
}

func toPrediction(value *types.Prediction) *Prediction {
	if value == nil {
		return nil
	}
	converted := &Prediction{Type: value.Type}
	if value.Content.Text != nil {
		converted.Content.Text = value.Content.Text
	}
	if value.Content.Parts != nil {
		parts := make([]TextContentPart, len(*value.Content.Parts))
		for i, part := range *value.Content.Parts {
			parts[i] = TextContentPart{Type: part.Type, Text: part.Text}
			if part.PromptCacheBreakpoint != nil {
				parts[i].PromptCacheBreakpoint = &PromptCacheBreakpoint{
					Mode: part.PromptCacheBreakpoint.Mode,
				}
			}
		}
		converted.Content.Parts = &parts
	}
	return converted
}

func toResponseFormat(value *types.ResponseFormat) *ResponseFormat {
	if value == nil {
		return nil
	}
	converted := &ResponseFormat{Type: ResponseFormatType(value.Type)}
	if value.JSONSchema != nil {
		converted.JSONSchema = &JSONSchemaFormat{
			Name:        value.JSONSchema.Name,
			Description: value.JSONSchema.Description,
			Schema:      value.JSONSchema.Schema,
			Strict:      value.JSONSchema.Strict,
		}
	}
	return converted
}

func toModalities(value *types.Modalities) *Modalities {
	if value == nil {
		return nil
	}
	converted := make(Modalities, len(*value))
	for i, modality := range *value {
		converted[i] = Modality(modality)
	}
	return &converted
}

func toAudioOutput(value *types.AudioOutput) *AudioOutput {
	if value == nil {
		return nil
	}
	return &AudioOutput{
		Format: AudioFormat(value.Format),
		Voice: AudioVoice{
			Name: value.Voice.Name,
			ID:   value.Voice.ID,
		},
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
			Index:        choice.Index,
			Message:      toResponseMessage(choice.Message),
			FinishReason: toFinishReason(choice.FinishReason),
			Logprobs:     toChatLogprobs(choice.Logprobs),
		})
	}

	return types.ChatCompletionResponse{
		ID:                resp.ID,
		Object:            resp.Object,
		Created:           resp.Created,
		Model:             resp.Model,
		Choices:           choices,
		Moderation:        toChatCompletionModeration(resp.Moderation),
		ServiceTier:       toResponseServiceTier(resp.ServiceTier),
		SystemFingerprint: resp.SystemFingerprint,
		Usage:             toCompletionUsage(resp.Usage),
	}
}

func toResponseServiceTier(value *ServiceTier) *types.ServiceTier {
	if value == nil {
		return nil
	}
	converted := types.ServiceTier(*value)
	return &converted
}

func toChatCompletionModeration(value *ChatCompletionModeration) *types.ChatCompletionModeration {
	if value == nil {
		return nil
	}
	return &types.ChatCompletionModeration{
		Input:  toModerationOutcome(value.Input),
		Output: toModerationOutcome(value.Output),
	}
}

func toModerationOutcome(value ModerationOutcome) types.ModerationOutcome {
	converted := types.ModerationOutcome{}
	if value.Results != nil {
		converted.Results = &types.ModerationResults{
			Model: value.Results.Model, Type: value.Results.Type,
			Results: toModerationResults(value.Results.Results),
		}
	}
	if value.Error != nil {
		converted.Error = &types.ModerationError{
			Code: value.Error.Code, Message: value.Error.Message, Type: value.Error.Type,
		}
	}
	return converted
}

func toModerationResults(values []ModerationResult) []types.ModerationResult {
	converted := make([]types.ModerationResult, len(values))
	for i, value := range values {
		converted[i] = types.ModerationResult{
			Categories: value.Categories, CategoryAppliedInputTypes: value.CategoryAppliedInputTypes,
			CategoryScores: value.CategoryScores, Flagged: value.Flagged, Model: value.Model, Type: value.Type,
		}
	}
	return converted
}

func toResponseMessage(value ResponseMessage) types.ResponseMessage {
	return types.ResponseMessage{
		Role:         value.Role,
		Content:      value.Content,
		Refusal:      value.Refusal,
		Annotations:  toResponseAnnotations(value.Annotations),
		Audio:        toResponseAudio(value.Audio),
		FunctionCall: toResponseFunctionCall(value.FunctionCall),
		ToolCalls:    toResponseToolCalls(value.ToolCalls),
	}
}

func toResponseAnnotations(value *[]Annotation) *[]types.Annotation {
	if value == nil {
		return nil
	}
	converted := make([]types.Annotation, len(*value))
	for i, annotation := range *value {
		converted[i] = types.Annotation{
			Type: annotation.Type,
			URLCitation: types.URLCitation{
				EndIndex: annotation.URLCitation.EndIndex, StartIndex: annotation.URLCitation.StartIndex,
				Title: annotation.URLCitation.Title, URL: annotation.URLCitation.URL,
			},
		}
	}
	return &converted
}

func toResponseAudio(value *ChatCompletionAudio) *types.ChatCompletionAudio {
	if value == nil {
		return nil
	}
	return &types.ChatCompletionAudio{
		ID: value.ID, Data: value.Data, ExpiresAt: value.ExpiresAt, Transcript: value.Transcript,
	}
}

func toResponseFunctionCall(value *MessageFunctionCall) *types.ResponseFunctionCall {
	if value == nil {
		return nil
	}
	return &types.ResponseFunctionCall{Arguments: value.Arguments, Name: value.Name}
}

func toResponseToolCalls(value *[]MessageToolCall) *[]types.ResponseToolCall {
	if value == nil {
		return nil
	}
	converted := make([]types.ResponseToolCall, len(*value))
	for i, call := range *value {
		converted[i] = types.ResponseToolCall{
			ID: call.ID, Type: call.Type, Function: toResponseFunctionCall(call.Function),
		}
		if call.Custom != nil {
			converted[i].Custom = &types.ResponseCustomCall{
				Input: call.Custom.Input, Name: call.Custom.Name,
			}
		}
	}
	return &converted
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
			FinishReason: toFinishReason(choice.FinishReason),
			Logprobs:     toChatLogprobs(choice.Logprobs),
		})
	}

	return types.StreamChunk{
		ID:      chunk.ID,
		Object:  chunk.Object,
		Created: chunk.Created,
		Model:   chunk.Model,
		Choices: choices,
		Usage:   toCompletionUsage(chunk.Usage),
	}
}

func toFinishReason(value *FinishReason) *types.FinishReason {
	if value == nil {
		return nil
	}
	converted := types.FinishReason(*value)
	return &converted
}

func toChatLogprobs(value *ChatLogprobs) *types.ChatLogprobs {
	if value == nil {
		return nil
	}
	return &types.ChatLogprobs{
		Content: toTokenLogprobs(value.Content),
		Refusal: toTokenLogprobs(value.Refusal),
	}
}

func toTokenLogprobs(values []TokenLogprob) []types.TokenLogprob {
	if values == nil {
		return nil
	}
	converted := make([]types.TokenLogprob, len(values))
	for i, value := range values {
		converted[i] = types.TokenLogprob{
			Token:       value.Token,
			Bytes:       value.Bytes,
			Logprob:     value.Logprob,
			TopLogprobs: toTopLogprobs(value.TopLogprobs),
		}
	}
	return converted
}

func toTopLogprobs(values []TopLogprob) []types.TopLogprob {
	if values == nil {
		return nil
	}
	converted := make([]types.TopLogprob, len(values))
	for i, value := range values {
		converted[i] = types.TopLogprob{
			Token: value.Token, Bytes: value.Bytes, Logprob: value.Logprob,
		}
	}
	return converted
}

func toCompletionUsage(value *CompletionUsage) *types.CompletionUsage {
	if value == nil {
		return nil
	}
	converted := &types.CompletionUsage{
		CompletionTokens: value.CompletionTokens,
		PromptTokens:     value.PromptTokens,
		TotalTokens:      value.TotalTokens,
	}
	if value.CompletionTokensDetails != nil {
		details := value.CompletionTokensDetails
		converted.CompletionTokensDetails = &types.CompletionTokensDetails{
			AcceptedPredictionTokens: details.AcceptedPredictionTokens,
			AudioTokens:              details.AudioTokens,
			ReasoningTokens:          details.ReasoningTokens,
			RejectedPredictionTokens: details.RejectedPredictionTokens,
		}
	}
	if value.PromptTokensDetails != nil {
		details := value.PromptTokensDetails
		converted.PromptTokensDetails = &types.PromptTokensDetails{
			AudioTokens:      details.AudioTokens,
			CacheWriteTokens: details.CacheWriteTokens,
			CachedTokens:     details.CachedTokens,
		}
	}
	return converted
}
