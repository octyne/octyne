package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoadReturnsDefaultOpenAIProvider(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-api-key")
	unsetEnv(t, "PORT")

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

func TestLoadRequiresOpenAIAPIKey(t *testing.T) {
	unsetEnv(t, "OPENAI_API_KEY")

	got, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want missing OPENAI_API_KEY error")
	}
	if !reflect.DeepEqual(got, Config{}) {
		t.Errorf("Load() config = %+v, want zero value", got)
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
