package openai

import (
	"bytes"
	"context"
	"encoding/json"
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

	body, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		a.config.BaseURL+"chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	return &types.ChatCompletionResponse{
		ID: "chatcmpl_openai",
	}, nil
}
