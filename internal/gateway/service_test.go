package gateway

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/registry"
	"github.com/octyne/octyne/internal/types"
)

func TestResolveRouteReturnsTypedRoutingErrors(t *testing.T) {
	tests := []struct {
		name             string
		model            string
		providerRegistry *providers.Registry
		modelRegistry    *registry.Registry
		wantKind         types.ErrorKind
		wantStatus       int
		wantParam        string
	}{
		{
			name:             "unknown model",
			model:            "not-registered",
			providerRegistry: providers.NewRegistry(),
			modelRegistry:    registry.NewRegistry(),
			wantKind:         types.ErrorKindNotFound,
			wantStatus:       http.StatusNotFound,
			wantParam:        "model",
		},
		{
			name:             "provider missing",
			model:            "openai/gpt-5-nano",
			providerRegistry: providers.NewRegistry(),
			modelRegistry:    registryWithModel("openai/gpt-5-nano", "openai", "gpt-5-nano"),
			wantKind:         types.ErrorKindInternal,
			wantStatus:       http.StatusInternalServerError,
		},
		{
			name:  "adapter missing",
			model: "openai/gpt-5-nano",
			providerRegistry: registryWithProvider(providers.New(
				providers.Config{Name: "openai"},
				nil,
			)),
			modelRegistry: registryWithModel("openai/gpt-5-nano", "openai", "gpt-5-nano"),
			wantKind:      types.ErrorKindInternal,
			wantStatus:    http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(
				tt.providerRegistry,
				tt.modelRegistry,
			).resolveRoute(tt.model)
			var apiErr *types.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("error = %T, want *types.APIError", err)
			}
			if apiErr.Kind != tt.wantKind || apiErr.HTTPStatus != tt.wantStatus {
				t.Errorf("error = %+v, want kind %q status %d", apiErr, tt.wantKind, tt.wantStatus)
			}
			if stringValue(apiErr.Param) != tt.wantParam {
				t.Errorf("Param = %q, want %q", stringValue(apiErr.Param), tt.wantParam)
			}
		})
	}
}

func TestChatUsesResolvedUpstreamModelID(t *testing.T) {
	adapter := &recordingAdapter{}
	service := New(
		registryWithProvider(providers.New(
			providers.Config{Name: "openai"},
			adapter,
		)),
		registryWithModel("openai/gpt-5-nano", "openai", "gpt-5-nano"),
	)
	request := types.ChatCompletionRequest{Model: "openai/gpt-5-nano"}

	if _, err := service.Chat(context.Background(), request); err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	if got := adapter.chatRequest.Model; got != "gpt-5-nano" {
		t.Errorf("adapter request model = %q, want %q", got, "gpt-5-nano")
	}
	if request.Model != "openai/gpt-5-nano" {
		t.Errorf("caller request model = %q, want public model %q", request.Model, "openai/gpt-5-nano")
	}
}

func TestStreamChatUsesResolvedUpstreamModelID(t *testing.T) {
	adapter := &recordingAdapter{}
	service := New(
		registryWithProvider(providers.New(
			providers.Config{Name: "openai"},
			adapter,
		)),
		registryWithModel("openai/gpt-5-nano", "openai", "gpt-5-nano"),
	)
	request := types.ChatCompletionRequest{Model: "openai/gpt-5-nano"}

	if _, err := service.StreamChat(context.Background(), request); err != nil {
		t.Fatalf("StreamChat() error = %v", err)
	}

	if got := adapter.streamRequest.Model; got != "gpt-5-nano" {
		t.Errorf("adapter request model = %q, want %q", got, "gpt-5-nano")
	}
	if request.Model != "openai/gpt-5-nano" {
		t.Errorf("caller request model = %q, want public model %q", request.Model, "openai/gpt-5-nano")
	}
}

type recordingAdapter struct {
	chatRequest   types.ChatCompletionRequest
	streamRequest types.ChatCompletionRequest
}

func (a *recordingAdapter) Chat(
	_ context.Context,
	req types.ChatCompletionRequest,
) (*types.ChatCompletionResponse, error) {
	a.chatRequest = req
	return &types.ChatCompletionResponse{}, nil
}

func (a *recordingAdapter) StreamChat(
	_ context.Context,
	req types.ChatCompletionRequest,
) (<-chan types.StreamChunk, error) {
	a.streamRequest = req
	chunks := make(chan types.StreamChunk)
	close(chunks)
	return chunks, nil
}

func registryWithModel(name string, provider string, modelID string) *registry.Registry {
	modelRegistry := registry.NewRegistry()
	modelRegistry.Register(name, registry.Model{
		Provider: provider,
		ModelID:  modelID,
	})
	return modelRegistry
}

func registryWithProvider(provider *providers.Provider) *providers.Registry {
	registry := providers.NewRegistry()
	registry.Register(provider.Name, provider)
	return registry
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
