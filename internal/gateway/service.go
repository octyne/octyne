package gateway

import (
	"context"
	"errors"

	"github.com/usekeel/keel/internal/providers"
	"github.com/usekeel/keel/internal/registry"
	"github.com/usekeel/keel/internal/types"
)

type Service struct {
	providers *providers.Registry
}

func New(providers *providers.Registry) *Service {
	return &Service{
		providers: providers,
	}
}

func (s *Service) Chat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (*types.ChatCompletionResponse, error) {

	model, ok := registry.Get(req.Model)
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

	return provider.Adapter.Chat(
		ctx,
		req,
	)
}
