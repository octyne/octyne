package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	providersEnv                          = "OCTYNE_PROVIDERS"
	defaultProviderNames                  = "openai"
	defaultOpenAIBaseURL                  = "https://api.openai.com/v1"
	defaultOpenAIModels                   = "gpt-4.1-mini,gpt-5-nano"
	defaultNonStreamingTimeout            = 600 * time.Second
	defaultStreamingResponseHeaderTimeout = 30 * time.Second
	clientAPIKeysEnv                      = "OCTYNE_API_KEYS"
	minimumClientAPIKeyLength             = 32
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

func providerEnvPrefix(name string) (string, error) {
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return "", fmt.Errorf(
			"provider name %q must not start or end with a hyphen",
			name,
		)
	}

	for _, character := range name {
		isLowercaseLetter := character >= 'a' && character <= 'z'
		isDigit := character >= '0' && character <= '9'

		if !isLowercaseLetter && !isDigit && character != '-' {
			return "", fmt.Errorf(
				"provider name %q must contain only lowercase letters, digits, and hyphens",
				name,
			)
		}
	}

	return strings.ToUpper(strings.ReplaceAll(name, "-", "_")), nil
}

func parseModels(providerName string, value string) ([]ModelConfig, error) {
	parts := strings.Split(value, ",")
	models := make([]ModelConfig, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		upstreamID := strings.TrimSpace(part)
		if upstreamID == "" {
			return nil, fmt.Errorf(
				"provider %q contains an empty model ID",
				providerName,
			)
		}
		if _, exists := seen[upstreamID]; exists {
			return nil, fmt.Errorf(
				"provider %q contains duplicate model ID %q",
				providerName,
				upstreamID,
			)
		}

		seen[upstreamID] = struct{}{}
		models = append(models, ModelConfig{
			PublicName: providerName + "/" + upstreamID,
			UpstreamID: upstreamID,
		})
	}

	return models, nil
}

func loadDuration(envName string, fallback time.Duration) (time.Duration, error) {
	value, exists := os.LookupEnv(envName)
	if !exists {
		return fallback, nil
	}

	duration, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf(
			"parse %s as duration: %w",
			envName,
			err,
		)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", envName)
	}

	return duration, nil
}

func normalizeBaseURL(envName string, value string) (string, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(value), "/")

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("parse %s as URL: %w", envName, err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("%s must use http or https", envName)
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("%s must include a host", envName)
	}
	if parsed.User != nil {
		return "", fmt.Errorf("%s must not contain credentials", envName)
	}
	if parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" {
		return "", fmt.Errorf("%s must not contain a query or fragment", envName)
	}

	return baseURL, nil
}

func loadProviderNames() ([]string, error) {
	value, exists := os.LookupEnv(providersEnv)
	if !exists {
		value = defaultProviderNames
	}

	parts := strings.Split(value, ",")
	names := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name == "" {
			return nil, errors.New("OCTYNE_PROVIDERS contains an empty provider name")
		}
		if _, err := providerEnvPrefix(name); err != nil {
			return nil, err
		}
		if _, exists := seen[name]; exists {
			return nil, fmt.Errorf(
				"OCTYNE_PROVIDERS contains duplicate provider %q",
				name,
			)
		}

		seen[name] = struct{}{}
		names = append(names, name)
	}

	return names, nil
}

func loadProvider(name string) (ProviderConfig, error) {
	prefix, err := providerEnvPrefix(name)
	if err != nil {
		return ProviderConfig{}, err
	}

	baseURLEnv := prefix + "_BASE_URL"
	baseURLValue, exists := os.LookupEnv(baseURLEnv)
	if !exists {
		if name != "openai" {
			return ProviderConfig{}, fmt.Errorf(
				"required environment variable %s is missing",
				baseURLEnv,
			)
		}
		baseURLValue = defaultOpenAIBaseURL
	}

	baseURL, err := normalizeBaseURL(baseURLEnv, baseURLValue)
	if err != nil {
		return ProviderConfig{}, err
	}

	apiKeyEnv := prefix + "_API_KEY"
	apiKey, exists := os.LookupEnv(apiKeyEnv)
	if !exists {
		return ProviderConfig{}, fmt.Errorf(
			"required environment variable %s is missing",
			apiKeyEnv,
		)
	}

	modelsEnv := prefix + "_MODELS"
	modelsValue, exists := os.LookupEnv(modelsEnv)
	if !exists {
		if name != "openai" {
			return ProviderConfig{}, fmt.Errorf(
				"required environment variable %s is missing",
				modelsEnv,
			)
		}
		modelsValue = defaultOpenAIModels
	}

	models, err := parseModels(name, modelsValue)
	if err != nil {
		return ProviderConfig{}, err
	}

	nonStreamingTimeout, err := loadDuration(
		prefix+"_NON_STREAMING_TIMEOUT",
		defaultNonStreamingTimeout,
	)
	if err != nil {
		return ProviderConfig{}, err
	}

	streamingResponseHeaderTimeout, err := loadDuration(
		prefix+"_STREAMING_RESPONSE_HEADER_TIMEOUT",
		defaultStreamingResponseHeaderTimeout,
	)
	if err != nil {
		return ProviderConfig{}, err
	}

	return ProviderConfig{
		Name:                           name,
		BaseURL:                        baseURL,
		APIKey:                         apiKey,
		NonStreamingTimeout:            nonStreamingTimeout,
		StreamingResponseHeaderTimeout: streamingResponseHeaderTimeout,
		Models:                         models,
	}, nil
}

func loadClientAPIKeys() ([]string, error) {
	value, exists := os.LookupEnv(clientAPIKeysEnv)
	if !exists {
		return nil, fmt.Errorf(
			"required environment variable %s is missing",
			clientAPIKeysEnv,
		)
	}

	parts := strings.Split(value, ",")
	keys := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for index, part := range parts {
		key := strings.TrimSpace(part)
		if key == "" {
			return nil, fmt.Errorf(
				"%s contains an empty key at position %d",
				clientAPIKeysEnv,
				index+1,
			)
		}
		if len(key) < minimumClientAPIKeyLength {
			return nil, fmt.Errorf(
				"%s key at position %d must contain at least %d characters",
				clientAPIKeysEnv,
				index+1,
				minimumClientAPIKeyLength,
			)
		}
		if _, exists := seen[key]; exists {
			return nil, fmt.Errorf(
				"%s contains a duplicate key at position %d",
				clientAPIKeysEnv,
				index+1,
			)
		}

		seen[key] = struct{}{}
		keys = append(keys, key)
	}

	return keys, nil
}

type Config struct {
	Port          string
	ClientAPIKeys []string
	Providers     []ProviderConfig
}

func Load() (Config, error) {

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "3000"
	}

	clientAPIKeys, err := loadClientAPIKeys()
	if err != nil {
		return Config{}, fmt.Errorf("load client API keys: %w", err)
	}

	providerNames, err := loadProviderNames()
	if err != nil {
		return Config{}, err
	}

	providerConfigs := make([]ProviderConfig, 0, len(providerNames))
	for _, name := range providerNames {
		providerConfig, err := loadProvider(name)
		if err != nil {
			return Config{}, fmt.Errorf(
				"load provider %q: %w",
				name,
				err,
			)
		}

		providerConfigs = append(providerConfigs, providerConfig)
	}

	return Config{
		Port:          port,
		ClientAPIKeys: clientAPIKeys,
		Providers:     providerConfigs,
	}, nil
}
