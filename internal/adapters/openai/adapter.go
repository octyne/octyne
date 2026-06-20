package openai

import (
	"context"
	"net/http"

	"github.com/usekeel/keel/internal/providers"
	"github.com/usekeel/keel/internal/types"
)

type Adapter struct {
	config providers.Config
	client *http.Client
}

func New(config providers.Config) *Adapter {
	return &Adapter{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (a *Adapter) Chat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (*types.ChatCompletionResponse, error) {
	openAIReq := toChatCompletionRequest(req)
	_ = openAIReq

	return &types.ChatCompletionResponse{
		ID: "chatcmpl_openai",
	}, nil
}
