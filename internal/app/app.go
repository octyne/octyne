package app

import (
	"log/slog"
	"time"

	"github.com/octyne/octyne/internal/adapters/openai"
	"github.com/octyne/octyne/internal/config"
	"github.com/octyne/octyne/internal/gateway"
	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/registry"
	"github.com/octyne/octyne/internal/server"
)

type App struct {
	Server *server.Server
}

func New(appConfig config.Config, logger *slog.Logger) *App {
	providerRegistry := providers.NewRegistry()

	cfg := providers.Config{
		Name:                           "openai",
		BaseURL:                        "https://api.openai.com/v1",
		APIKey:                         appConfig.OpenAIAPIKey,
		NonStreamingTimeout:            600 * time.Second,
		StreamingResponseHeaderTimeout: 30 * time.Second,
	}

	providerRegistry.Register(
		"openai",
		providers.New(
			cfg,
			openai.New(cfg),
		),
	)

	modelRegistry := registry.NewRegistry()

	modelRegistry.Register(
		"openai/gpt-4.1-mini",
		registry.Model{
			Provider: "openai",
			ModelID:  "gpt-4.1-mini",
		},
	)

	modelRegistry.Register(
		"openai/gpt-5-nano",
		registry.Model{
			Provider: "openai",
			ModelID:  "gpt-5-nano",
		},
	)

	gatewayService := gateway.New(
		providerRegistry,
		modelRegistry,
	)

	httpServer := server.New(
		":"+appConfig.Port,
		logger,
		gatewayService,
		modelRegistry,
	)

	return &App{
		Server: httpServer,
	}
}
