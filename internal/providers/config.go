package providers

import "time"

type Config struct {
	Name    string
	BaseURL string
	APIKey  string
	Timeout time.Duration
}
