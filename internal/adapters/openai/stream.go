package openai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/octyne/octyne/internal/types"
)

func readChatCompletionStream(
	ctx context.Context,
	body io.ReadCloser,
) <-chan types.StreamChunk {
	chunks := make(chan types.StreamChunk)

	go func() {
		defer close(chunks)
		defer body.Close()

		scanner := bufio.NewScanner(body)
		scanner.Buffer(
			make([]byte, 64*1024),
			1024*1024,
		)

		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data:") {
				continue
			}

			data := strings.TrimSpace(
				strings.TrimPrefix(line, "data:"),
			)

			if data == "" {
				continue
			}

			if data == "[DONE]" {
				return
			}

			var openAIChunk ChatCompletionChunk

			if err := json.Unmarshal(
				[]byte(data),
				&openAIChunk,
			); err != nil {
				select {
				case chunks <- types.StreamChunk{
					Error: fmt.Errorf(
						"decode openai stream chunk: %w",
						err,
					),
				}:
				case <-ctx.Done():
				}

				return
			}

			select {
			case chunks <- toStreamChunk(openAIChunk):
			case <-ctx.Done():
				return
			}
		}

		if err := scanner.Err(); err != nil && ctx.Err() == nil {
			select {
			case chunks <- types.StreamChunk{
				Error: fmt.Errorf(
					"read openai stream: %w",
					err,
				),
			}:
			case <-ctx.Done():
			}
		}
	}()

	return chunks
}
