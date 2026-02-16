package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Domain      string
	AccessToken string
	Retry       RetryConfig
}

type RetryConfig struct {
	MaxRetries    int
	InitialDelay  time.Duration
	BackoffFactor float64
}

func Load() (*Config, error) {
	domain := os.Getenv("TP_DOMAIN")
	token := os.Getenv("TP_ACCESS_TOKEN")
	if domain == "" || token == "" {
		return nil, fmt.Errorf("TP_DOMAIN and TP_ACCESS_TOKEN environment variables are required")
	}
	return &Config{
		Domain:      domain,
		AccessToken: token,
		Retry: RetryConfig{
			MaxRetries:    3,
			InitialDelay:  1 * time.Second,
			BackoffFactor: 2.0,
		},
	}, nil
}
