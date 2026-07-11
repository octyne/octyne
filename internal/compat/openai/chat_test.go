package openai

import (
	"encoding/json"
	"testing"
)

func TestStopSequencesAcceptsStringAndArray(t *testing.T) {
	tests := []struct {
		name string
		json string
		want []string
	}{
		{name: "single", json: `"done"`, want: []string{"done"}},
		{name: "array", json: `["done","stop"]`, want: []string{"done", "stop"}},
		{name: "empty string", json: `""`, want: []string{""}},
		{name: "empty array", json: `[]`, want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got StopSequences
			if err := json.Unmarshal([]byte(tt.json), &got); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("len(StopSequences) = %d, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("StopSequences[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestAudioVoiceAcceptsNameAndCustomID(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantName string
		wantID   string
	}{
		{name: "name", json: `"alloy"`, wantName: "alloy"},
		{name: "custom", json: `{"id":"voice_123"}`, wantID: "voice_123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got AudioVoice
			if err := json.Unmarshal([]byte(tt.json), &got); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if tt.wantName != "" && (got.Name == nil || *got.Name != tt.wantName) {
				t.Errorf("Name = %v, want %q", got.Name, tt.wantName)
			}
			if tt.wantID != "" && (got.ID == nil || *got.ID != tt.wantID) {
				t.Errorf("ID = %v, want %q", got.ID, tt.wantID)
			}
		})
	}
}

func TestPredictionContentAcceptsStringAndParts(t *testing.T) {
	var text PredictionContent
	if err := json.Unmarshal([]byte(`"hello"`), &text); err != nil {
		t.Fatalf("unmarshal text: %v", err)
	}
	if text.Text == nil || *text.Text != "hello" || text.Parts != nil {
		t.Errorf("text prediction = %+v", text)
	}

	var parts PredictionContent
	if err := json.Unmarshal([]byte(`[{"type":"text","text":"hello"}]`), &parts); err != nil {
		t.Fatalf("unmarshal parts: %v", err)
	}
	if parts.Parts == nil || len(*parts.Parts) != 1 || (*parts.Parts)[0].Text != "hello" {
		t.Errorf("parts prediction = %+v", parts)
	}
}

func TestToolChoiceAcceptsModeAndObject(t *testing.T) {
	var mode ToolChoice
	if err := json.Unmarshal([]byte(`"auto"`), &mode); err != nil {
		t.Fatalf("unmarshal mode: %v", err)
	}
	if mode.Mode == nil || *mode.Mode != "auto" {
		t.Errorf("Mode = %v, want auto", mode.Mode)
	}

	var named ToolChoice
	if err := json.Unmarshal([]byte(`{"type":"function","function":{"name":"weather"}}`), &named); err != nil {
		t.Fatalf("unmarshal named choice: %v", err)
	}
	if named.Function == nil || named.Function.Name != "weather" {
		t.Errorf("Function = %+v, want weather", named.Function)
	}
}

func TestLegacyFunctionCallAcceptsModeAndName(t *testing.T) {
	var mode LegacyFunctionCall
	if err := json.Unmarshal([]byte(`"none"`), &mode); err != nil {
		t.Fatalf("unmarshal mode: %v", err)
	}
	if mode.Mode == nil || *mode.Mode != "none" {
		t.Errorf("Mode = %v, want none", mode.Mode)
	}

	var named LegacyFunctionCall
	if err := json.Unmarshal([]byte(`{"name":"legacy"}`), &named); err != nil {
		t.Fatalf("unmarshal name: %v", err)
	}
	if named.Name == nil || *named.Name != "legacy" {
		t.Errorf("Name = %v, want legacy", named.Name)
	}
}
