package providers

import "github.com/usekeel/keel/internal/adapters"

type Provider struct {
	Name    string
	Config  Config
	Adapter adapters.Adapter
}
