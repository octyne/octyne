package gateway

import (
	"context"
	"net/http"

	"github.com/octyne/octyne/internal/adapters"
	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/registry"
	"github.com/octyne/octyne/internal/types"
)

type Service struct {
	providers *providers.Registry
	models    *registry.Registry
}

func New(providers *providers.Registry, models *registry.Registry) *Service {
	return &Service{
		providers: providers,
		models:    models,
	}
}

func (s *Service) resolveAdapter(modelName string) (adapters.Adapter, error) {
	model, ok := s.models.Get(modelName)
	if !ok {
		param := "model"
		code := "model_not_found"
		return nil, &types.APIError{
			Kind:       types.ErrorKindNotFound,
			Message:    "The requested model does not exist.",
			Param:      &param,
			Code:       &code,
			HTTPStatus: http.StatusNotFound,
		}
	}

	provider, ok := s.providers.Get(model.Provider)
	if !ok {
		return nil, &types.APIError{
			Kind:       types.ErrorKindInternal,
			Message:    "The model provider is not available.",
			HTTPStatus: http.StatusInternalServerError,
		}
	}

	if provider.Adapter == nil {
		return nil, &types.APIError{
			Kind:       types.ErrorKindInternal,
			Message:    "The model provider is not configured.",
			HTTPStatus: http.StatusInternalServerError,
		}
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
