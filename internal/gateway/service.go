package gateway

import "github.com/usekeel/keel/internal/types"

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Chat(req types.ChatCompletionRequest) (*types.ChatCompletionResponse, error) {
	return &types.ChatCompletionResponse{
		ID: "chatcmpl_id",
	}, nil
}
