package adapters

import (
	"context"

	"github.com/usekeel/keel/internal/types"
)

type Adapter interface {
	Chat(
		ctx context.Context,
		req types.ChatCompletionRequest,
	) (*types.ChatCompletionResponse, error)
}
