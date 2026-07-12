package app

import (
	"log/slog"

	"github.com/octyne/octyne/internal/adapters/openai"
	"github.com/octyne/octyne/internal/auth"
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
	clientKeyVerifier := auth.NewStaticKeyVerifier(appConfig.ClientAPIKeys)
	providerRegistry := providers.NewRegistry()
	modelRegistry := registry.NewRegistry()

	for _, providerConfig := range appConfig.Providers {
		runtimeConfig := providers.Config{
			Name:                           providerConfig.Name,
			BaseURL:                        providerConfig.BaseURL,
			APIKey:                         providerConfig.APIKey,
			NonStreamingTimeout:            providerConfig.NonStreamingTimeout,
			StreamingResponseHeaderTimeout: providerConfig.StreamingResponseHeaderTimeout,
		}

		providerRegistry.Register(
			providerConfig.Name,
			providers.New(
				runtimeConfig,
				openai.New(runtimeConfig),
			),
		)
		for _, modelConfig := range providerConfig.Models {
			modelRegistry.Register(
				modelConfig.PublicName,
				registry.Model{
					Provider: providerConfig.Name,
					ModelID:  modelConfig.UpstreamID,
				},
			)
		}
	}

	gatewayService := gateway.New(
		providerRegistry,
		modelRegistry,
	)

	httpServer := server.New(
		":"+appConfig.Port,
		logger,
		gatewayService,
		modelRegistry,
		clientKeyVerifier,
	)

	return &App{
		Server: httpServer,
	}
}
