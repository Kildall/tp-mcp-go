package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tp-mcp-go/internal/domain/entity"
)

// FetchMetadata fetches the TP API metadata
func (c *httpClient) FetchMetadata(ctx context.Context) (any, error) {
	url := fmt.Sprintf("%s/Index/meta", c.baseURL)
	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	var result any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetValidEntityTypes returns cached entity types or fetches from API
func (c *httpClient) GetValidEntityTypes(ctx context.Context) ([]string, error) {
	// Check cache first
	c.cacheMu.RLock()
	if c.cachedTypes != nil && time.Now().Before(c.cacheExpiry) {
		types := make([]string, len(c.cachedTypes))
		copy(types, c.cachedTypes)
		c.cacheMu.RUnlock()
		return types, nil
	}
	c.cacheMu.RUnlock()

	// Try to fetch from API metadata
	metadata, err := c.FetchMetadata(ctx)
	if err != nil {
		// Fallback to static types
		return staticEntityTypes(), nil
	}

	// Extract entity type names from metadata
	types := extractEntityTypes(metadata)
	if len(types) == 0 {
		return staticEntityTypes(), nil
	}

	// Update cache
	c.cacheMu.Lock()
	c.cachedTypes = types
	c.cacheExpiry = time.Now().Add(1 * time.Hour)
	c.cacheMu.Unlock()

	return types, nil
}

// InitializeCache pre-populates the entity type cache
func (c *httpClient) InitializeCache(ctx context.Context) error {
	_, err := c.GetValidEntityTypes(ctx)
	return err
}

// extractEntityTypes pulls entity type names from the metadata response
func extractEntityTypes(metadata any) []string {
	// Metadata is expected to be a map with entity type info
	metaMap, ok := metadata.(map[string]any)
	if !ok {
		return nil
	}

	var types []string
	for key := range metaMap {
		types = append(types, key)
	}
	return types
}

// staticEntityTypes returns the hardcoded list of valid entity types
func staticEntityTypes() []string {
	types := make([]string, len(entity.ValidTypes))
	for i, t := range entity.ValidTypes {
		types[i] = string(t)
	}
	return types
}
