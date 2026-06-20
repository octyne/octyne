package providers

import (
	"context"

	"github.com/usekeel/keel/internal/types"
)

type Provider interface {
	Chat(
		ctx context.Context,
		req types.ChatCompletionRequest,
	) (*types.ChatCompletionResponse, error)
}
