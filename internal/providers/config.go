package providers

import "time"

type Config struct {
	Name                           string
	BaseURL                        string
	APIKey                         string
	NonStreamingTimeout            time.Duration
	StreamingResponseHeaderTimeout time.Duration
}
