package gateway

import (
	"context"
	"errors"

	"github.com/octyne/octyne/internal/adapters"
	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/registry"
	"github.com/octyne/octyne/internal/types"
)

type Service struct {
	providers *providers.Registry
}

func New(providers *providers.Registry) *Service {
	return &Service{
		providers: providers,
	}
}

func (s *Service) resolveAdapter(modelName string) (adapters.Adapter, error) {
	model, ok := registry.Get(modelName)
	if !ok {
		return nil, errors.New("unknown model")
	}

	provider, ok := s.providers.Get(model.Provider)
	if !ok {
		return nil, errors.New("provider not found")
	}

	if provider.Adapter == nil {
		return nil, errors.New("provider adapter not configured")
	}

	return provider.Adapter, nil
}

func (s *Service) Chat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (*types.ChatCompletionResponse, error) {

	adapter, err := s.resolveAdapter(req.Model)
	if err != nil {
		return nil, err
	}

	return adapter.Chat(
		ctx,
		req,
	)
}

func (s *Service) StreamChat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (<-chan types.StreamChunk, error) {

	adapter, err := s.resolveAdapter(req.Model)
	if err != nil {
		return nil, err
	}

	return adapter.StreamChat(
		ctx,
		req,
	)
}
