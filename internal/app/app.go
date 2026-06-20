package app

import (
	"time"

	"github.com/usekeel/keel/internal/adapters/openai"
	"github.com/usekeel/keel/internal/gateway"
	"github.com/usekeel/keel/internal/providers"
	"github.com/usekeel/keel/internal/server"
)

type App struct {
	Server *server.Server
}

func New() *App {
	providerRegistry := providers.NewRegistry()

	cfg := providers.Config{
		Name:    "openai",
		BaseURL: "https://api.openai.com/v1",
		Timeout: 30 * time.Second,
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
