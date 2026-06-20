package gateway

import (
	"errors"

	"github.com/usekeel/keel/internal/registry"
	"github.com/usekeel/keel/internal/types"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Chat(req types.ChatCompletionRequest) (*types.ChatCompletionResponse, error) {

	_, ok := registry.Get(req.Model)
	if !ok {
		return nil, errors.New("unknown model")
	}

	return &types.ChatCompletionResponse{
		ID: "chatcmpl_id",
	}, nil
}
