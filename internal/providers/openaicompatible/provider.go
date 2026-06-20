package openaicompatible

import (
	"context"

	"github.com/usekeel/keel/internal/types"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Chat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (*types.ChatCompletionResponse, error) {
	return &types.ChatCompletionResponse{
		ID: "chatcmple_openai",
	}, nil
}
