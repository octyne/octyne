package types

type StreamChunk struct {
	Content string
	Error   error
	Done    bool
}
