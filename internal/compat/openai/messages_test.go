package openai

import (
	"encoding/json"
	"testing"
)

func TestMessagesDecodeRoleSpecificShapes(t *testing.T) {
	var request ChatCompletionRequest
	err := json.Unmarshal([]byte(`{
		"model":"openai/gpt-5-nano",
		"messages":[
			{"role":"developer","content":[{"type":"text","text":"rules","prompt_cache_breakpoint":{"mode":"explicit"}}]},
			{"role":"system","content":"system"},
			{"role":"user","content":[{"type":"image_url","image_url":{"url":"data:image/png;base64,AA==","detail":"low"}},{"type":"input_audio","input_audio":{"data":"AA==","format":"wav"}},{"type":"file","file":{"file_id":"file_1"}}]},
			{"role":"assistant","content":null,"audio":{"id":"audio_1"},"refusal":"","tool_calls":[{"id":"call_1","type":"function","function":{"name":"weather","arguments":"{}"}},{"id":"call_2","type":"custom","custom":{"name":"shell","input":"pwd"}}]},
			{"role":"tool","content":[{"type":"text","text":"sunny"}],"tool_call_id":"call_1"},
			{"role":"function","content":"legacy result","name":"legacy"}
		]
	}`), &request)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(request.Messages) != 6 {
		t.Fatalf("len(Messages) = %d, want 6", len(request.Messages))
	}
	if request.Messages[0].Developer == nil || request.Messages[0].Developer.Content.Parts == nil {
		t.Error("developer message not decoded")
	}
	if request.Messages[2].User == nil || request.Messages[2].User.Content.Parts == nil || len(*request.Messages[2].User.Content.Parts) != 3 {
		t.Error("multimodal user message not decoded")
	}
	assistant := request.Messages[3].Assistant
	if assistant == nil || !assistant.ContentNull || assistant.ToolCalls == nil || len(*assistant.ToolCalls) != 2 {
		t.Errorf("assistant = %+v", assistant)
	}
	if request.Messages[4].Tool == nil || request.Messages[4].Tool.ToolCallID != "call_1" {
		t.Error("tool message not decoded")
	}
	if request.Messages[5].Function == nil || request.Messages[5].Function.Name != "legacy" {
		t.Error("function message not decoded")
	}
}
