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
