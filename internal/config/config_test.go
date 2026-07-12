package config

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestLoadReturnsDefaultOpenAIProvider(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-api-key")
	unsetEnv(t, "PORT")
	unsetEnv(t, "OCTYNE_PROVIDERS")
	unsetEnv(t, "OPENAI_BASE_URL")
	unsetEnv(t, "OPENAI_MODELS")
	unsetEnv(t, "OPENAI_NON_STREAMING_TIMEOUT")
	unsetEnv(t, "OPENAI_STREAMING_RESPONSE_HEADER_TIMEOUT")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := Config{
		Port: "3000",
		Providers: []ProviderConfig{
			{
				Name:                           "openai",
				BaseURL:                        "https://api.openai.com/v1",
				APIKey:                         "test-api-key",
				NonStreamingTimeout:            600 * time.Second,
				StreamingResponseHeaderTimeout: 30 * time.Second,
				Models: []ModelConfig{
					{
						PublicName: "openai/gpt-4.1-mini",
						UpstreamID: "gpt-4.1-mini",
					},
					{
						PublicName: "openai/gpt-5-nano",
						UpstreamID: "gpt-5-nano",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load() = %+v, want %+v", got, want)
	}
}

func TestLoadReturnsConfiguredOpenAIAndOllamaProviders(t *testing.T) {
	t.Setenv("PORT", "4321")
	t.Setenv("OCTYNE_PROVIDERS", " openai, ollama ")
	t.Setenv("OPENAI_BASE_URL", "https://openai.example/v1/")
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("OPENAI_MODELS", "gpt-5-nano")
	t.Setenv("OPENAI_NON_STREAMING_TIMEOUT", "2m")
	t.Setenv("OPENAI_STREAMING_RESPONSE_HEADER_TIMEOUT", "15s")
	t.Setenv("OLLAMA_BASE_URL", "http://localhost:11434/v1/")
	t.Setenv("OLLAMA_API_KEY", "")
	t.Setenv("OLLAMA_MODELS", "llama3.2, qwen3:8b")
	t.Setenv("OLLAMA_NON_STREAMING_TIMEOUT", "5m")
	t.Setenv("OLLAMA_STREAMING_RESPONSE_HEADER_TIMEOUT", "45s")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := Config{
		Port: "4321",
		Providers: []ProviderConfig{
			{
				Name:                           "openai",
				BaseURL:                        "https://openai.example/v1",
				APIKey:                         "openai-key",
				NonStreamingTimeout:            2 * time.Minute,
				StreamingResponseHeaderTimeout: 15 * time.Second,
				Models: []ModelConfig{
					{PublicName: "openai/gpt-5-nano", UpstreamID: "gpt-5-nano"},
				},
			},
			{
				Name:                           "ollama",
				BaseURL:                        "http://localhost:11434/v1",
				APIKey:                         "",
				NonStreamingTimeout:            5 * time.Minute,
				StreamingResponseHeaderTimeout: 45 * time.Second,
				Models: []ModelConfig{
					{PublicName: "ollama/llama3.2", UpstreamID: "llama3.2"},
					{PublicName: "ollama/qwen3:8b", UpstreamID: "qwen3:8b"},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load() = %+v, want %+v", got, want)
	}
}

func TestProviderEnvPrefixConvertsKebabCase(t *testing.T) {
	got, err := providerEnvPrefix("azure-openai")
	if err != nil {
		t.Fatalf("providerEnvPrefix() error = %v", err)
	}
	if got != "AZURE_OPENAI" {
		t.Errorf("providerEnvPrefix() = %q, want %q", got, "AZURE_OPENAI")
	}
}

func TestLoadRequiresOpenAIAPIKey(t *testing.T) {
	unsetEnv(t, "OCTYNE_PROVIDERS")
	unsetEnv(t, "OPENAI_API_KEY")

	got, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want missing OPENAI_API_KEY error")
	}
	if !reflect.DeepEqual(got, Config{}) {
		t.Errorf("Load() config = %+v, want zero value", got)
	}
}

func TestLoadRejectsInvalidProviderLists(t *testing.T) {
	tests := []struct {
		name      string
		providers string
		wantError string
	}{
		{
			name:      "empty name",
			providers: "openai,,ollama",
			wantError: "empty provider name",
		},
		{
			name:      "duplicate name",
			providers: "openai,openai",
			wantError: `duplicate provider "openai"`,
		},
		{
			name:      "uppercase name",
			providers: "OpenAI",
			wantError: "only lowercase letters",
		},
		{
			name:      "leading hyphen",
			providers: "-openai",
			wantError: "must not start or end with a hyphen",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OCTYNE_PROVIDERS", test.providers)

			got, err := Load()
			if err == nil {
				t.Fatal("Load() error = nil, want validation error")
			}
			if !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("Load() error = %q, want substring %q", err, test.wantError)
			}
			if !reflect.DeepEqual(got, Config{}) {
				t.Errorf("Load() config = %+v, want zero value", got)
			}
		})
	}
}

func TestLoadRequiresExplicitNonOpenAIProviderValues(t *testing.T) {
	tests := []struct {
		name      string
		missing   string
		wantError string
	}{
		{name: "base URL", missing: "OLLAMA_BASE_URL", wantError: "OLLAMA_BASE_URL is missing"},
		{name: "API key declaration", missing: "OLLAMA_API_KEY", wantError: "OLLAMA_API_KEY is missing"},
		{name: "models", missing: "OLLAMA_MODELS", wantError: "OLLAMA_MODELS is missing"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OCTYNE_PROVIDERS", "ollama")
			t.Setenv("OLLAMA_BASE_URL", "http://localhost:11434/v1")
			t.Setenv("OLLAMA_API_KEY", "")
			t.Setenv("OLLAMA_MODELS", "llama3.2")
			unsetEnv(t, test.missing)

			got, err := Load()
			if err == nil {
				t.Fatal("Load() error = nil, want missing environment variable error")
			}
			if !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("Load() error = %q, want substring %q", err, test.wantError)
			}
			if !reflect.DeepEqual(got, Config{}) {
				t.Errorf("Load() config = %+v, want zero value", got)
			}
		})
	}
}

func TestLoadRejectsInvalidProviderValues(t *testing.T) {
	tests := []struct {
		name      string
		envName   string
		value     string
		wantError string
	}{
		{
			name:      "base URL scheme",
			envName:   "OLLAMA_BASE_URL",
			value:     "file:///tmp/ollama",
			wantError: "must use http or https",
		},
		{
			name:      "base URL credentials",
			envName:   "OLLAMA_BASE_URL",
			value:     "http://user:password@localhost:11434/v1",
			wantError: "must not contain credentials",
		},
		{
			name:      "duplicate model",
			envName:   "OLLAMA_MODELS",
			value:     "llama3.2,llama3.2",
			wantError: `duplicate model ID "llama3.2"`,
		},
		{
			name:      "invalid timeout",
			envName:   "OLLAMA_NON_STREAMING_TIMEOUT",
			value:     "soon",
			wantError: "parse OLLAMA_NON_STREAMING_TIMEOUT as duration",
		},
		{
			name:      "zero timeout",
			envName:   "OLLAMA_STREAMING_RESPONSE_HEADER_TIMEOUT",
			value:     "0s",
			wantError: "must be greater than zero",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OCTYNE_PROVIDERS", "ollama")
			t.Setenv("OLLAMA_BASE_URL", "http://localhost:11434/v1")
			t.Setenv("OLLAMA_API_KEY", "")
			t.Setenv("OLLAMA_MODELS", "llama3.2")
			t.Setenv(test.envName, test.value)

			got, err := Load()
			if err == nil {
				t.Fatal("Load() error = nil, want validation error")
			}
			if !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("Load() error = %q, want substring %q", err, test.wantError)
			}
			if !reflect.DeepEqual(got, Config{}) {
				t.Errorf("Load() config = %+v, want zero value", got)
			}
		})
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	value, existed := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}

	t.Cleanup(func() {
		if existed {
			if err := os.Setenv(key, value); err != nil {
				t.Errorf("restore %s: %v", key, err)
			}
			return
		}

		if err := os.Unsetenv(key); err != nil {
			t.Errorf("unset %s during cleanup: %v", key, err)
		}
	})
}
