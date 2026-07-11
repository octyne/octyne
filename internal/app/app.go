package app

import (
	"time"

	"github.com/octyne/octyne/internal/adapters/openai"
	"github.com/octyne/octyne/internal/config"
	"github.com/octyne/octyne/internal/gateway"
	"github.com/octyne/octyne/internal/providers"
	"github.com/octyne/octyne/internal/server"
)

type App struct {
	Server *server.Server
}

func New(appConfig config.Config) *App {
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

	gatewayService := gateway.New(
		providerRegistry,
	)

	httpServer := server.New(
		gatewayService,
	)

	return &App{
		Server: httpServer,
	}
}
