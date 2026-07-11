package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/octyne/octyne/internal/adapters/openai"
	"github.com/octyne/octyne/internal/gateway"
	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/types"
)

func TestChatHandlerStreamsOpenAICompatibleSSE(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/chat/completions" {
				t.Errorf(
					"path = %q, want /chat/completions",
					r.URL.Path,
				)
			}

			if got := r.Header.Get("Accept"); got != "text/event-stream" {
				t.Errorf(
					"Accept = %q, want text/event-stream",
					got,
				)
			}

			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read request body: %v", err)
				return
			}

			if strings.Contains(string(requestBody), `"temperature"`) {
				t.Errorf(
					"request body unexpectedly includes temperature: %s",
					requestBody,
				)
			}

			var upstreamRequest struct {
				TopP                 *float64           `json:"top_p"`
				FrequencyPenalty     *float64           `json:"frequency_penalty"`
				PresencePenalty      *float64           `json:"presence_penalty"`
				MaxCompletionTokens  *int               `json:"max_completion_tokens"`
				N                    *int               `json:"n"`
				Logprobs             *bool              `json:"logprobs"`
				TopLogprobs          *int               `json:"top_logprobs"`
				ReasoningEffort      *string            `json:"reasoning_effort"`
				Verbosity            *string            `json:"verbosity"`
				Seed                 *int64             `json:"seed"`
				Store                *bool              `json:"store"`
				ParallelToolCalls    *bool              `json:"parallel_tool_calls"`
				SafetyIdentifier     *string            `json:"safety_identifier"`
				PromptCacheKey       *string            `json:"prompt_cache_key"`
				MaxTokens            *int               `json:"max_tokens"`
				User                 *string            `json:"user"`
				PromptCacheRetention *string            `json:"prompt_cache_retention"`
				Metadata             *map[string]string `json:"metadata"`
				ServiceTier          *string            `json:"service_tier"`
				PromptCacheOptions   *struct {
					Mode *string `json:"mode"`
					TTL  *string `json:"ttl"`
				} `json:"prompt_cache_options"`
				Stop          *[]string           `json:"stop"`
				LogitBias     *map[string]float64 `json:"logit_bias"`
				StreamOptions *struct {
					IncludeUsage       *bool `json:"include_usage"`
					IncludeObfuscation *bool `json:"include_obfuscation"`
				} `json:"stream_options"`
			}

			if err := json.Unmarshal(requestBody, &upstreamRequest); err != nil {
				t.Errorf("decode upstream request: %v", err)
				return
			}

			if upstreamRequest.TopP == nil {
				t.Error("TopP = nil, want explicit zero")
			} else if *upstreamRequest.TopP != 0 {
				t.Errorf(
					"TopP = %v, want 0",
					*upstreamRequest.TopP,
				)
			}

			if upstreamRequest.FrequencyPenalty != nil {
				t.Errorf(
					"FrequencyPenalty = %v, want nil",
					*upstreamRequest.FrequencyPenalty,
				)
			}

			if upstreamRequest.PresencePenalty != nil {
				t.Errorf(
					"PresencePenalty = %v, want nil",
					*upstreamRequest.PresencePenalty,
				)
			}

			if upstreamRequest.MaxCompletionTokens == nil {
				t.Error("MaxCompletionTokens = nil, want 128")
			} else if *upstreamRequest.MaxCompletionTokens != 128 {
				t.Errorf(
					"MaxCompletionTokens = %v, want 128",
					*upstreamRequest.MaxCompletionTokens,
				)
			}

			if upstreamRequest.N == nil {
				t.Error("N = nil, want 2")
			} else if *upstreamRequest.N != 2 {
				t.Errorf("N = %v, want 2", *upstreamRequest.N)
			}

			if upstreamRequest.Logprobs == nil {
				t.Error("Logprobs = nil, want true")
			} else if !*upstreamRequest.Logprobs {
				t.Error("Logprobs = false, want true")
			}

			if upstreamRequest.TopLogprobs == nil {
				t.Error("TopLogprobs = nil, want explicit zero")
			} else if *upstreamRequest.TopLogprobs != 0 {
				t.Errorf(
					"TopLogprobs = %d, want 0",
					*upstreamRequest.TopLogprobs,
				)
			}

			if upstreamRequest.ReasoningEffort == nil ||
				*upstreamRequest.ReasoningEffort != "high" {
				t.Errorf(
					"ReasoningEffort = %v, want high",
					upstreamRequest.ReasoningEffort,
				)
			}

			if upstreamRequest.Verbosity != nil {
				t.Errorf("Verbosity = %v, want nil", *upstreamRequest.Verbosity)
			}

			if upstreamRequest.Seed == nil || *upstreamRequest.Seed != 0 {
				t.Errorf("Seed = %v, want explicit zero", upstreamRequest.Seed)
			}
			if upstreamRequest.Store == nil || *upstreamRequest.Store {
				t.Errorf("Store = %v, want explicit false", upstreamRequest.Store)
			}
			if upstreamRequest.ParallelToolCalls == nil ||
				*upstreamRequest.ParallelToolCalls {
				t.Errorf(
					"ParallelToolCalls = %v, want explicit false",
					upstreamRequest.ParallelToolCalls,
				)
			}
			if upstreamRequest.SafetyIdentifier == nil ||
				*upstreamRequest.SafetyIdentifier != "" {
				t.Errorf("SafetyIdentifier = %v, want empty", upstreamRequest.SafetyIdentifier)
			}
			if upstreamRequest.PromptCacheKey == nil ||
				*upstreamRequest.PromptCacheKey != "" {
				t.Errorf("PromptCacheKey = %v, want empty", upstreamRequest.PromptCacheKey)
			}
			if upstreamRequest.MaxTokens == nil || *upstreamRequest.MaxTokens != 0 {
				t.Errorf("MaxTokens = %v, want explicit zero", upstreamRequest.MaxTokens)
			}
			if upstreamRequest.User == nil || *upstreamRequest.User != "" {
				t.Errorf("User = %v, want empty", upstreamRequest.User)
			}
			if upstreamRequest.PromptCacheRetention == nil ||
				*upstreamRequest.PromptCacheRetention != "24h" {
				t.Errorf(
					"PromptCacheRetention = %v, want 24h",
					upstreamRequest.PromptCacheRetention,
				)
			}
			if upstreamRequest.Metadata == nil || len(*upstreamRequest.Metadata) != 0 {
				t.Errorf("Metadata = %v, want explicit empty object", upstreamRequest.Metadata)
			}
			if upstreamRequest.ServiceTier == nil || *upstreamRequest.ServiceTier != "flex" {
				t.Errorf("ServiceTier = %v, want flex", upstreamRequest.ServiceTier)
			}
			if upstreamRequest.PromptCacheOptions == nil ||
				upstreamRequest.PromptCacheOptions.Mode == nil ||
				*upstreamRequest.PromptCacheOptions.Mode != "explicit" ||
				upstreamRequest.PromptCacheOptions.TTL == nil ||
				*upstreamRequest.PromptCacheOptions.TTL != "30m" {
				t.Errorf(
					"PromptCacheOptions = %+v, want explicit/30m",
					upstreamRequest.PromptCacheOptions,
				)
			}
			if upstreamRequest.Stop == nil || len(*upstreamRequest.Stop) != 1 ||
				(*upstreamRequest.Stop)[0] != "END" {
				t.Errorf("Stop = %v, want [END]", upstreamRequest.Stop)
			}
			if upstreamRequest.LogitBias == nil || len(*upstreamRequest.LogitBias) != 0 {
				t.Errorf("LogitBias = %v, want explicit empty object", upstreamRequest.LogitBias)
			}
			if upstreamRequest.StreamOptions == nil ||
				upstreamRequest.StreamOptions.IncludeUsage == nil ||
				*upstreamRequest.StreamOptions.IncludeUsage ||
				upstreamRequest.StreamOptions.IncludeObfuscation == nil ||
				*upstreamRequest.StreamOptions.IncludeObfuscation {
				t.Errorf(
					"StreamOptions = %+v, want explicit false values",
					upstreamRequest.StreamOptions,
				)
			}

			if !strings.Contains(
				string(requestBody),
				`"stream":true`,
			) {
				t.Errorf(
					"request body does not enable streaming: %s",
					requestBody,
				)
			}

			w.Header().Set("Content-Type", "text/event-stream")

			_, err = io.WriteString(w, `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":123,"model":"gpt-5-nano","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null,"logprobs":null}]}

data: [DONE]

`)
			if err != nil {
				t.Errorf("write upstream stream: %v", err)
			}
		},
	))

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(
			`{"model":"gpt-5-nano","messages":[{"role":"user","content":"Hello"}],"stream":true,"top_p":0,"max_completion_tokens":128,"n":2,"logprobs":true,"top_logprobs":0,"reasoning_effort":"high","seed":0,"store":false,"parallel_tool_calls":false,"safety_identifier":"","prompt_cache_key":"","max_tokens":0,"user":"","prompt_cache_retention":"24h","metadata":{},"service_tier":"flex","prompt_cache_options":{"mode":"explicit","ttl":"30m"},"stop":"END","logit_bias":{},"stream_options":{"include_usage":false,"include_obfuscation":false}}`,
		),
	)

	recorder := httptest.NewRecorder()
	server.mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d",
			recorder.Code,
			http.StatusOK,
		)
	}

	if got := recorder.Header().Get("Content-Type"); got != "text/event-stream" {
		t.Errorf(
			"Content-Type = %q, want text/event-stream",
			got,
		)
	}

	wantBody := `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":123,"model":"gpt-5-nano","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null,"logprobs":null}]}

data: [DONE]

`

	if got := recorder.Body.String(); got != wantBody {
		t.Errorf(
			"response body:\n%s\nwant:\n%s",
			got,
			wantBody,
		)
	}

	if !recorder.Flushed {
		t.Error("response was not flushed")
	}
}

func TestChatHandlerReturnsOpenAICompatibleJSON(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read request body: %v", err)
				return
			}

			var upstreamRequest struct {
				Temperature         *float64 `json:"temperature"`
				TopP                *float64 `json:"top_p"`
				FrequencyPenalty    *float64 `json:"frequency_penalty"`
				PresencePenalty     *float64 `json:"presence_penalty"`
				MaxCompletionTokens *int     `json:"max_completion_tokens"`
				N                   *int     `json:"n"`
				Logprobs            *bool    `json:"logprobs"`
				TopLogprobs         *int     `json:"top_logprobs"`
				ReasoningEffort     *string  `json:"reasoning_effort"`
				Verbosity           *string  `json:"verbosity"`
				Seed                *int64   `json:"seed"`
				Store               *bool    `json:"store"`
				ParallelToolCalls   *bool    `json:"parallel_tool_calls"`
				SafetyIdentifier    *string  `json:"safety_identifier"`
				PromptCacheKey      *string  `json:"prompt_cache_key"`
			}

			if err := json.Unmarshal(requestBody, &upstreamRequest); err != nil {
				t.Errorf("decode upstream request: %v", err)
				return
			}

			if upstreamRequest.Temperature == nil {
				t.Error("Temperature = nil, want explicit zero")
			} else if *upstreamRequest.Temperature != 0 {
				t.Errorf(
					"Temperature = %v, want 0",
					*upstreamRequest.Temperature,
				)
			}

			if upstreamRequest.TopP != nil {
				t.Errorf(
					"TopP = %v, want nil",
					*upstreamRequest.TopP,
				)
			}

			if upstreamRequest.FrequencyPenalty == nil {
				t.Error("FrequencyPenalty = nil, want explicit zero")
			} else if *upstreamRequest.FrequencyPenalty != 0 {
				t.Errorf(
					"FrequencyPenalty = %v, want 0",
					*upstreamRequest.FrequencyPenalty,
				)
			}

			if upstreamRequest.PresencePenalty == nil {
				t.Error("PresencePenalty = nil, want explicit zero")
			} else if *upstreamRequest.PresencePenalty != 0 {
				t.Errorf(
					"PresencePenalty = %v, want 0",
					*upstreamRequest.PresencePenalty,
				)
			}

			if upstreamRequest.MaxCompletionTokens != nil {
				t.Errorf(
					"MaxCompletionTokens = %v, want nil",
					*upstreamRequest.MaxCompletionTokens,
				)
			}

			if upstreamRequest.N != nil {
				t.Errorf("N = %v, want nil", *upstreamRequest.N)
			}

			if upstreamRequest.Logprobs == nil {
				t.Error("Logprobs = nil, want explicit false")
			} else if *upstreamRequest.Logprobs {
				t.Error("Logprobs = true, want false")
			}

			if upstreamRequest.TopLogprobs != nil {
				t.Errorf(
					"TopLogprobs = %d, want nil",
					*upstreamRequest.TopLogprobs,
				)
			}

			if upstreamRequest.ReasoningEffort != nil {
				t.Errorf(
					"ReasoningEffort = %v, want nil",
					*upstreamRequest.ReasoningEffort,
				)
			}

			if upstreamRequest.Verbosity == nil ||
				*upstreamRequest.Verbosity != "medium" {
				t.Errorf(
					"Verbosity = %v, want medium",
					upstreamRequest.Verbosity,
				)
			}

			if upstreamRequest.Seed != nil || upstreamRequest.Store != nil ||
				upstreamRequest.ParallelToolCalls != nil ||
				upstreamRequest.SafetyIdentifier != nil ||
				upstreamRequest.PromptCacheKey != nil {
				t.Errorf("omitted scalar controls reached upstream: %+v", upstreamRequest)
			}

			if strings.Contains(string(requestBody), `"stream":true`) {
				t.Errorf(
					"non-streaming request enables streaming: %s",
					requestBody,
				)
			}

			w.Header().Set("Content-Type", "application/json")

			_, err = io.WriteString(w, `{
				"id":"chatcmpl-456",
				"object":"chat.completion",
				"created":456,
				"model":"gpt-5-nano",
				"choices":[{
					"index":0,
					"message":{"role":"assistant","content":"Hello"},
					"finish_reason":"stop",
					"logprobs":null
				}],
				"usage":{"prompt_tokens":5,"completion_tokens":1,"total_tokens":6}
			}`)
			if err != nil {
				t.Errorf("write upstream response: %v", err)
			}
		},
	))

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(
			`{"model":"gpt-5-nano","messages":[{"role":"user","content":"Hello"}],"temperature":0,"frequency_penalty":0,"presence_penalty":0,"logprobs":false,"verbosity":"medium"}`,
		),
	)

	recorder := httptest.NewRecorder()
	server.mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d: %s",
			recorder.Code,
			http.StatusOK,
			recorder.Body.String(),
		)
	}

	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	var response types.ChatCompletionResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.ID != "chatcmpl-456" ||
		response.Object != "chat.completion" ||
		response.Created != 456 ||
		response.Model != "gpt-5-nano" {
		t.Errorf("unexpected response metadata: %+v", response)
	}

	if len(response.Choices) != 1 {
		t.Fatalf("len(Choices) = %d, want 1", len(response.Choices))
	}

	choice := response.Choices[0]
	if choice.Index != 0 ||
		choice.Message.Role != "assistant" ||
		choice.Message.Content != "Hello" ||
		choice.FinishReason == nil ||
		*choice.FinishReason != "stop" {
		t.Errorf("unexpected choice: %+v", choice)
	}

	if response.Usage == nil {
		t.Error("Usage = nil, want OpenAI usage object")
	}
}

func newTestServer(t *testing.T, upstreamHandler http.Handler) *Server {
	t.Helper()

	upstream := httptest.NewServer(upstreamHandler)
	t.Cleanup(upstream.Close)

	config := providers.Config{
		Name:                           "openai",
		BaseURL:                        upstream.URL,
		NonStreamingTimeout:            time.Second,
		StreamingResponseHeaderTimeout: time.Second,
	}

	providerRegistry := providers.NewRegistry()
	providerRegistry.Register(
		"openai",
		providers.New(
			config,
			openai.New(config),
		),
	)

	return New(gateway.New(providerRegistry))
}

func TestChatHandlerRejectsEmptyMessages(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t.Error("upstream was called for an invalid request")
			http.Error(w, "unexpected upstream call", http.StatusInternalServerError)
		},
	))

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/chat/completions",
		strings.NewReader(
			`{"model":"gpt-5-nano","messages":[]}`,
		),
	)

	recorder := httptest.NewRecorder()
	server.mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(
			"status = %d, want %d",
			recorder.Code,
			http.StatusBadRequest,
		)
	}

	if got := recorder.Body.String(); got != "messages are required\n" {
		t.Errorf(
			"body = %q, want %q",
			got,
			"messages are required\n",
		)
	}
}
