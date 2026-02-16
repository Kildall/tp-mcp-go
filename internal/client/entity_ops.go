package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"tp-mcp-go/internal/domain/entity"
)

// GetEntity retrieves an entity by type and ID
func (c *httpClient) GetEntity(ctx context.Context, entityType entity.Type, id int, include []string) (map[string]any, error) {
	url := fmt.Sprintf("%s/%d", c.buildURL(entityType), id)
	if len(include) > 0 {
		url += fmt.Sprintf("?include=[%s]", strings.Join(include, ","))
	}
	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateEntity creates a new entity
func (c *httpClient) CreateEntity(ctx context.Context, entityType entity.Type, data map[string]any) (map[string]any, error) {
	url := c.buildURL(entityType)
	respData, err := c.doPost(ctx, url, data)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateEntity updates an existing entity
func (c *httpClient) UpdateEntity(ctx context.Context, entityType entity.Type, id int, data map[string]any) (map[string]any, error) {
	url := fmt.Sprintf("%s/%d", c.buildURL(entityType), id)
	respData, err := c.doPost(ctx, url, data)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, err
	}
	return result, nil
}
