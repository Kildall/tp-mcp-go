package auth

import (
	"net/http"
	"testing"
)

func TestApplyAuth_AddsAccessToken(t *testing.T) {
	strategy := NewAccessTokenStrategy("test-token-123")
	req, err := http.NewRequest("GET", "https://api.example.com/endpoint", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	strategy.ApplyAuth(req)

	token := req.URL.Query().Get("access_token")
	if token != "test-token-123" {
		t.Errorf("expected access_token=test-token-123, got %s", token)
	}
}

func TestApplyAuth_PreservesExistingQueryParameters(t *testing.T) {
	strategy := NewAccessTokenStrategy("my-token")
	req, err := http.NewRequest("GET", "https://api.example.com/endpoint?foo=bar&baz=qux", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	strategy.ApplyAuth(req)

	query := req.URL.Query()
	if query.Get("access_token") != "my-token" {
		t.Errorf("expected access_token=my-token, got %s", query.Get("access_token"))
	}
	if query.Get("foo") != "bar" {
		t.Errorf("expected foo=bar, got %s", query.Get("foo"))
	}
	if query.Get("baz") != "qux" {
		t.Errorf("expected baz=qux, got %s", query.Get("baz"))
	}
}

func TestApplyAuth_EmptyExistingQuery(t *testing.T) {
	strategy := NewAccessTokenStrategy("empty-test-token")
	req, err := http.NewRequest("GET", "https://api.example.com/endpoint", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	strategy.ApplyAuth(req)

	query := req.URL.Query()
	if len(query) != 1 {
		t.Errorf("expected exactly 1 query parameter, got %d", len(query))
	}
	if query.Get("access_token") != "empty-test-token" {
		t.Errorf("expected access_token=empty-test-token, got %s", query.Get("access_token"))
	}
}
