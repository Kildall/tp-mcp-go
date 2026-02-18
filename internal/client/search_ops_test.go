package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"tp-mcp-go/internal/client/auth"
	"tp-mcp-go/internal/config"
	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/domain/query"
)

// emptyAPIResponse returns a minimal valid TP API response JSON
func emptyAPIResponse() []byte {
	resp := map[string]any{
		"Items": []any{},
		"Next":  "",
	}
	b, _ := json.Marshal(resp)
	return b
}

// newTestClient creates an httpClient pointing at the given server URL with no retries.
func newTestClient(serverURL string) *httpClient {
	return &httpClient{
		baseURL:    serverURL + "/api/v1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
		auth:       auth.NewAccessTokenStrategy("test-token"),
		retryConfig: config.RetryConfig{
			MaxRetries:    0,
			InitialDelay:  0,
			BackoffFactor: 1.0,
		},
		token: "test-token",
	}
}

func TestSearchEntities_OrderByAscending(t *testing.T) {
	var capturedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(emptyAPIResponse())
	}))
	defer server.Close()

	c := newTestClient(server.URL)

	req := query.SearchRequest{
		EntityType:   entity.Type("UserStory"),
		OrderByField: "CreateDate",
		OrderByDesc:  false,
	}

	_, err := c.SearchEntities(context.Background(), req)
	if err != nil {
		t.Fatalf("SearchEntities returned unexpected error: %v", err)
	}

	if !strings.Contains(capturedURL, "orderBy=CreateDate") {
		t.Errorf("expected URL to contain 'orderBy=CreateDate', got: %s", capturedURL)
	}
	if strings.Contains(capturedURL, "orderByDesc") {
		t.Errorf("expected URL NOT to contain 'orderByDesc', got: %s", capturedURL)
	}
}

func TestSearchEntities_OrderByDescending(t *testing.T) {
	var capturedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(emptyAPIResponse())
	}))
	defer server.Close()

	c := newTestClient(server.URL)

	req := query.SearchRequest{
		EntityType:   entity.Type("UserStory"),
		OrderByField: "CreateDate",
		OrderByDesc:  true,
	}

	_, err := c.SearchEntities(context.Background(), req)
	if err != nil {
		t.Fatalf("SearchEntities returned unexpected error: %v", err)
	}

	if !strings.Contains(capturedURL, "orderByDesc=CreateDate") {
		t.Errorf("expected URL to contain 'orderByDesc=CreateDate', got: %s", capturedURL)
	}
	// The plain orderBy param (without Desc suffix) should not appear independently.
	// After stripping the "orderByDesc" occurrences, "orderBy=" should not be present.
	stripped := strings.ReplaceAll(capturedURL, "orderByDesc", "")
	if strings.Contains(stripped, "orderBy=") {
		t.Errorf("expected URL NOT to contain plain 'orderBy=' when descending, got: %s", capturedURL)
	}
}

func TestSearchEntities_NoOrderBy(t *testing.T) {
	var capturedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(emptyAPIResponse())
	}))
	defer server.Close()

	c := newTestClient(server.URL)

	req := query.SearchRequest{
		EntityType:   entity.Type("UserStory"),
		OrderByField: "",
	}

	_, err := c.SearchEntities(context.Background(), req)
	if err != nil {
		t.Fatalf("SearchEntities returned unexpected error: %v", err)
	}

	if strings.Contains(capturedURL, "orderBy") {
		t.Errorf("expected URL NOT to contain 'orderBy' (any form), got: %s", capturedURL)
	}
}
