package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/types"
)

type Adapter struct {
	config       providers.Config
	client       *http.Client
	streamClient *http.Client
}

func New(config providers.Config) *Adapter {
	nonStreamingTransport := http.DefaultTransport.(*http.Transport).Clone()
	streamingTransport := http.DefaultTransport.(*http.Transport).Clone()
	streamingTransport.ResponseHeaderTimeout = config.StreamingResponseHeaderTimeout

	return &Adapter{
		config: config,
		client: &http.Client{
			Transport: nonStreamingTransport,
			Timeout:   config.NonStreamingTimeout,
		},
		streamClient: &http.Client{
			Transport: streamingTransport,
		},
	}
}

func (a *Adapter) newChatRequest(
	ctx context.Context,
	req types.ChatCompletionRequest,
	stream bool,
) (*http.Request, error) {
	openAIReq := toChatCompletionRequest(req, stream)

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

	return httpReq, nil
}

func (a *Adapter) doChatRequest(
	client *http.Client,
	httpReq *http.Request,
) (*http.Response, error) {
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		responseBody, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"openai returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return resp, nil
}

func (a *Adapter) Chat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (*types.ChatCompletionResponse, error) {

	httpReq, err := a.newChatRequest(ctx, req, false)
	if err != nil {
		return nil, err
	}

	resp, err := a.doChatRequest(a.client, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openAIResp ChatCompletionResponse

	if err := json.Unmarshal(
		responseBody,
		&openAIResp,
	); err != nil {
		return nil, err
	}

	response := toChatCompletionResponse(
		openAIResp,
	)

	return &response, nil
}

func (a *Adapter) StreamChat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (<-chan types.StreamChunk, error) {
	httpReq, err := a.newChatRequest(ctx, req, true)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := a.doChatRequest(a.streamClient, httpReq)
	if err != nil {
		return nil, err
	}

	return readChatCompletionStream(
		ctx,
		resp.Body,
	), nil
}
