package config

import (
	"errors"
	"os"
)

type Config struct {
	Port         string
	OpenAIAPIKey string
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
		Port:         port,
		OpenAIAPIKey: openaiAPIKey,
	}, nil
}
