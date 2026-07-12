package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	openaicompat "github.com/octyne/octyne/internal/compat/openai"
	"github.com/octyne/octyne/internal/registry"
)

func TestModelsHandlerReturnsRegisteredModels(t *testing.T) {
	modelRegistry := registry.NewRegistry()
	modelRegistry.Register("openai/gpt-5-nano", registry.Model{
		Provider: "openai",
		ModelID:  "gpt-5-nano",
	})
	modelRegistry.Register("openai/gpt-4.1-mini", registry.Model{
		Provider: "openai",
		ModelID:  "gpt-4.1-mini",
	})
	server := New(nil, modelRegistry)

	request := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	recorder := httptest.NewRecorder()
	server.mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if got := recorder.Header().Get("x-request-id"); !strings.HasPrefix(got, "req_") {
		t.Errorf("x-request-id = %q, want req_ prefix", got)
	}

	var response openaicompat.ModelList
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Object != "list" {
		t.Errorf("Object = %q, want list", response.Object)
	}

	want := []openaicompat.Model{
		{
			ID:      "openai/gpt-4.1-mini",
			Object:  "model",
			Created: 0,
			OwnedBy: "openai",
		},
		{
			ID:      "openai/gpt-5-nano",
			Object:  "model",
			Created: 0,
			OwnedBy: "openai",
		},
	}
	if len(response.Data) != len(want) {
		t.Fatalf("len(Data) = %d, want %d", len(response.Data), len(want))
	}
	for i := range want {
		if response.Data[i] != want[i] {
			t.Errorf("Data[%d] = %+v, want %+v", i, response.Data[i], want[i])
		}
	}
}

func TestModelsHandlerReturnsEmptyList(t *testing.T) {
	server := New(nil, registry.NewRegistry())

	request := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	recorder := httptest.NewRecorder()
	server.mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Body.String(); got != "{\"object\":\"list\",\"data\":[]}\n" {
		t.Errorf("body = %q, want empty OpenAI model list", got)
	}
}

func TestModelsHandlerRejectsWrongMethod(t *testing.T) {
	server := New(nil, registry.NewRegistry())

	request := httptest.NewRequest(http.MethodPost, "/v1/models", nil)
	recorder := httptest.NewRecorder()
	server.mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
	if got := recorder.Header().Get("x-request-id"); !strings.HasPrefix(got, "req_") {
		t.Errorf("x-request-id = %q, want req_ prefix", got)
	}

	var response openAIErrorEnvelope
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Error.Code == nil || *response.Error.Code != "method_not_allowed" {
		t.Errorf("error code = %v, want method_not_allowed", response.Error.Code)
	}
}
