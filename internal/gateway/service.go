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

type route struct {
	adapter adapters.Adapter
	modelID string
}

func New(providers *providers.Registry, models *registry.Registry) *Service {
	return &Service{
		providers: providers,
		models:    models,
	}
}

func (s *Service) resolveRoute(modelName string) (route, error) {
	model, ok := s.models.Get(modelName)
	if !ok {
		param := "model"
		code := "model_not_found"
		return route{}, &types.APIError{
			Kind:       types.ErrorKindNotFound,
			Message:    "The requested model does not exist.",
			Param:      &param,
			Code:       &code,
			HTTPStatus: http.StatusNotFound,
		}
	}

	provider, ok := s.providers.Get(model.Provider)
	if !ok {
		return route{}, &types.APIError{
			Kind:       types.ErrorKindInternal,
			Message:    "The model provider is not available.",
			HTTPStatus: http.StatusInternalServerError,
		}
	}

	if provider.Adapter == nil {
		return route{}, &types.APIError{
			Kind:       types.ErrorKindInternal,
			Message:    "The model provider is not configured.",
			HTTPStatus: http.StatusInternalServerError,
		}
	}

	return route{
		adapter: provider.Adapter,
		modelID: model.ModelID,
	}, nil
}

func (s *Service) Chat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (*types.ChatCompletionResponse, error) {

	resolved, err := s.resolveRoute(req.Model)
	if err != nil {
		return nil, err
	}

	req.Model = resolved.modelID

	return resolved.adapter.Chat(
		ctx,
		req,
	)
}

func (s *Service) StreamChat(
	ctx context.Context,
	req types.ChatCompletionRequest,
) (<-chan types.StreamChunk, error) {

	resolved, err := s.resolveRoute(req.Model)
	if err != nil {
		return nil, err
	}

	req.Model = resolved.modelID

	return resolved.adapter.StreamChat(
		ctx,
		req,
	)
}
