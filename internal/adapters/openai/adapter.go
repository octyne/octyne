package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		a.config.BaseURL+"/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if a.config.APIKey != "" {
		httpReq.Header.Set(
			"Authorization",
			"Bearer "+a.config.APIKey,
		)
	}

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"openai returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// var openAIResp ChatCompletionResponse

	if err := json.Unmarshal(
		responseBody,
		&openAIReq,
	); err != nil {
		return nil, err
	}

	return &types.ChatCompletionResponse{
		ID: "chatcmpl_openai",
	}, nil
}
