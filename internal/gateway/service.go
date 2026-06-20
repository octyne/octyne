package gateway

import "github.com/usekeel/keel/internal/types"

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Chat(req types.ChatCompletionRequest) error {
	return nil
}
