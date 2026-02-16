package config

import (
	"testing"
	"time"
)

func TestLoad_BothEnvVarsSet(t *testing.T) {
	t.Setenv("TP_DOMAIN", "test.tpondemand.com")
	t.Setenv("TP_ACCESS_TOKEN", "test-token-123")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Domain != "test.tpondemand.com" {
		t.Errorf("Domain = %q, want %q", cfg.Domain, "test.tpondemand.com")
	}
	if cfg.AccessToken != "test-token-123" {
		t.Errorf("AccessToken = %q, want %q", cfg.AccessToken, "test-token-123")
	}
	if cfg.Retry.MaxRetries != 3 {
		t.Errorf("Retry.MaxRetries = %d, want %d", cfg.Retry.MaxRetries, 3)
	}
	if cfg.Retry.InitialDelay != 1*time.Second {
		t.Errorf("Retry.InitialDelay = %v, want %v", cfg.Retry.InitialDelay, 1*time.Second)
	}
	if cfg.Retry.BackoffFactor != 2.0 {
		t.Errorf("Retry.BackoffFactor = %f, want %f", cfg.Retry.BackoffFactor, 2.0)
	}
}

func TestLoad_MissingTPDomain(t *testing.T) {
	t.Setenv("TP_ACCESS_TOKEN", "test-token-123")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Load() expected error, got nil")
	}
	if cfg != nil {
		t.Errorf("Load() expected nil config, got %v", cfg)
	}
}

func TestLoad_MissingTPAccessToken(t *testing.T) {
	t.Setenv("TP_DOMAIN", "test.tpondemand.com")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Load() expected error, got nil")
	}
	if cfg != nil {
		t.Errorf("Load() expected nil config, got %v", cfg)
	}
}

func TestLoad_BothMissing(t *testing.T) {
	cfg, err := Load()
	if err == nil {
		t.Fatal("Load() expected error, got nil")
	}
	if cfg != nil {
		t.Errorf("Load() expected nil config, got %v", cfg)
	}
}

func TestLoad_DefaultRetryValues(t *testing.T) {
	t.Setenv("TP_DOMAIN", "test.tpondemand.com")
	t.Setenv("TP_ACCESS_TOKEN", "test-token-123")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Retry.MaxRetries != 3 {
		t.Errorf("Retry.MaxRetries = %d, want 3", cfg.Retry.MaxRetries)
	}
	if cfg.Retry.InitialDelay != 1*time.Second {
		t.Errorf("Retry.InitialDelay = %v, want 1s", cfg.Retry.InitialDelay)
	}
	if cfg.Retry.BackoffFactor != 2.0 {
		t.Errorf("Retry.BackoffFactor = %f, want 2.0", cfg.Retry.BackoffFactor)
	}
}
