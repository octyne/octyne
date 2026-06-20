package openaicompatible

import (
	"context"
	"net/http"

	"github.com/usekeel/keel/internal/providers"
	"github.com/usekeel/keel/internal/types"
)

type Provider struct {
	config providers.Config
	client *http.Client
}

func New(config providers.Config) *Provider {
	return &Provider{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (p *Provider) Chat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (*types.ChatCompletionResponse, error) {
	return &types.ChatCompletionResponse{
		ID: "chatcmple_openai",
	}, nil
}
