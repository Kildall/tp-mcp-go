package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"tp-mcp-go/internal/client/auth"
	"tp-mcp-go/internal/config"
	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/domain/errors"
)

type httpClient struct {
	baseURL     string
	httpClient  *http.Client
	auth        auth.Strategy
	retryConfig config.RetryConfig
	token       string

	// Entity type cache
	cacheMu     sync.RWMutex
	cachedTypes []string
	cacheExpiry time.Time
}

// NewHTTPClient creates a new Client implementation
func NewHTTPClient(cfg *config.Config, authStrategy auth.Strategy) Client {
	return &httpClient{
		baseURL:     fmt.Sprintf("https://%s/api/v1", cfg.Domain),
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		auth:        authStrategy,
		retryConfig: cfg.Retry,
		token:       cfg.AccessToken,
	}
}

// buildURL constructs the API URL for an entity type
func (c *httpClient) buildURL(entityType entity.Type) string {
	return fmt.Sprintf("%s/%s", c.baseURL, entity.Pluralize(entityType))
}

// doRequest executes an HTTP request with auth and retry
func (c *httpClient) doRequest(ctx context.Context, method, url string, body any) ([]byte, error) {
	return executeWithRetry(func() ([]byte, error) {
		var bodyReader io.Reader
		if body != nil {
			jsonBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(jsonBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, err
		}

		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		// Add format=json query param
		q := req.URL.Query()
		q.Set("format", "json")
		req.URL.RawQuery = q.Encode()

		// Apply auth (adds access_token param)
		c.auth.ApplyAuth(req)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, &errors.APIError{
				StatusCode: resp.StatusCode,
				Message:    errors.MaskToken(string(respBody), c.token),
				Context:    fmt.Sprintf("%s %s", method, url),
			}
		}

		return respBody, nil
	}, c.retryConfig.MaxRetries, c.retryConfig.InitialDelay, c.retryConfig.BackoffFactor)
}

// doGet performs a GET request
func (c *httpClient) doGet(ctx context.Context, url string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, url, nil)
}

// doPost performs a POST request with JSON body
func (c *httpClient) doPost(ctx context.Context, url string, body any) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPost, url, body)
}
