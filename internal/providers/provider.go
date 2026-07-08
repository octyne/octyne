package providers

import "github.com/octyne/octyne/internal/adapters"

type Provider struct {
	Name    string
	Config  Config
	Adapter adapters.Adapter
}

func New(
	config Config,
	adapter adapters.Adapter,
) *Provider {
	return &Provider{
		Name:    config.Name,
		Config:  config,
		Adapter: adapter,
	}

}
