package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"

	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/requestid"
	"github.com/octyne/octyne/internal/types"
)

const maxErrorBodyBytes = 64 * 1024

type upstreamErrorEnvelope struct {
	Error *upstreamError `json:"error"`
}

type upstreamError struct {
	Message string  `json:"message"`
	Param   *string `json:"param"`
	Code    *string `json:"code"`
}

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

	if id := requestid.FromContext(ctx); id != "" {
		httpReq.Header.Set("X-Client-Request-Id", id)
	}

	return httpReq, nil
}

func (a *Adapter) doChatRequest(
	client *http.Client,
	httpReq *http.Request,
) (*http.Response, error) {
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, toTransportError(httpReq.Context(), err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		responseBody, readErr := io.ReadAll(io.LimitReader(
			resp.Body,
			maxErrorBodyBytes+1,
		))

		return nil, toAPIError(resp, responseBody, readErr)
	}

	return resp, nil
}

func toTransportError(ctx context.Context, err error) *types.APIError {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return &types.APIError{
			Kind:       types.ErrorKindTimeout,
			Message:    "The upstream provider request timed out.",
			HTTPStatus: http.StatusGatewayTimeout,
			Cause:      err,
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return &types.APIError{
			Kind:       types.ErrorKindTimeout,
			Message:    "The upstream provider request timed out.",
			HTTPStatus: http.StatusGatewayTimeout,
			Cause:      err,
		}
	}

	return &types.APIError{
		Kind:       types.ErrorKindUnavailable,
		Message:    "The upstream provider is unavailable.",
		HTTPStatus: http.StatusBadGateway,
		Cause:      err,
	}
}

func toAPIError(resp *http.Response, body []byte, readErr error) *types.APIError {
	status := resp.StatusCode
	kind := types.ErrorKindInternal
	message := "The upstream provider failed to process the request."
	if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
		kind = types.ErrorKindInvalidRequest
		message = "The upstream provider rejected the request."
	}

	switch status {
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		kind = types.ErrorKindInvalidRequest
		message = "The upstream provider rejected the request."
	case http.StatusUnauthorized:
		kind = types.ErrorKindAuthentication
		message = "The upstream provider rejected authentication."
	case http.StatusForbidden:
		kind = types.ErrorKindPermission
		message = "The upstream provider denied the request."
	case http.StatusNotFound:
		kind = types.ErrorKindNotFound
		message = "The requested upstream resource was not found."
	case http.StatusTooManyRequests:
		kind = types.ErrorKindRateLimit
		message = "Rate limit exceeded."
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		kind = types.ErrorKindTimeout
		message = "The upstream provider request timed out."
	case http.StatusServiceUnavailable:
		kind = types.ErrorKindUnavailable
		message = "The upstream provider is unavailable."
	}

	apiErr := &types.APIError{
		Kind:              kind,
		Message:           message,
		HTTPStatus:        status,
		ProviderRequestID: resp.Header.Get("x-request-id"),
		Cause:             readErr,
	}

	if len(body) <= maxErrorBodyBytes {
		var envelope upstreamErrorEnvelope
		if err := json.Unmarshal(body, &envelope); err == nil &&
			envelope.Error != nil {
			if status < http.StatusInternalServerError &&
				envelope.Error.Message != "" {
				apiErr.Message = envelope.Error.Message
			}
			apiErr.Param = envelope.Error.Param
			apiErr.Code = envelope.Error.Code
		}
	}

	return apiErr
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
