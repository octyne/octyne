package types

type MessageContent struct {
	Text  *string
	Parts *[]ContentPart
}
type ImageURL struct {
	URL    string
	Detail *string
}
type InputAudio struct {
	Data   string
	Format string
}
type FileInput struct {
	FileData *string
	FileID   *string
	Filename *string
}
type ContentPart struct {
	Type                  string
	Text                  *string
	ImageURL              *ImageURL
	InputAudio            *InputAudio
	File                  *FileInput
	Refusal               *string
	PromptCacheBreakpoint *PromptCacheBreakpoint
}
type AudioReference struct{ ID string }
type MessageFunctionCall struct {
	Arguments string
	Name      string
}
type MessageCustomCall struct {
	Input string
	Name  string
}
type MessageToolCall struct {
	ID       string
	Type     string
	Function *MessageFunctionCall
	Custom   *MessageCustomCall
}
type ChatMessage struct {
	Role         string
	Content      *MessageContent
	Name         *string
	Audio        *AudioReference
	Refusal      *string
	ToolCalls    *[]MessageToolCall
	FunctionCall *MessageFunctionCall
	ToolCallID   *string
	ContentNull  bool
}
