package config

import (
	"errors"
	"os"
	"time"
)

type ProviderConfig struct {
	Name                           string
	BaseURL                        string
	APIKey                         string
	NonStreamingTimeout            time.Duration
	StreamingResponseHeaderTimeout time.Duration
	Models                         []ModelConfig
}

type ModelConfig struct {
	PublicName string
	UpstreamID string
}

type Config struct {
	Port      string
	Providers []ProviderConfig
}

func Load() (Config, error) {

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "3000"
	}

	openaiAPIKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		return Config{}, errors.New("required environment variable OPENAI_API_KEY is missing")
	}

	return Config{
		Port: port,
		Providers: []ProviderConfig{
			{
				Name:                           "openai",
				BaseURL:                        "https://api.openai.com/v1",
				APIKey:                         openaiAPIKey,
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
	}, nil
}
