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

func (s *Service) Chat(req types.ChatCompletionRequest) (*types.ChatCompletionResponse, error) {

	model, ok := registry.Get(req.Model)
	if !ok {
		return nil, errors.New("unknown model")
	}

	provider, ok := s.providers.Get(model.Provider)
	if !ok {
		return nil, errors.New("provider not found")
	}

	return provider.Chat(
		context.Background(),
		req,
	)
}
