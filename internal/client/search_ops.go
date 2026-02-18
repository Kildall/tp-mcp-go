package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/domain/query"
)

func (c *httpClient) SearchEntities(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error) {
	// If cursor is provided, validate SSRF and use cursor URL directly
	if req.Cursor != "" {
		if err := validateURL(req.Cursor, c.baseURL); err != nil {
			return nil, err
		}
		data, err := c.doGet(ctx, req.Cursor)
		if err != nil {
			return nil, err
		}
		return c.parseSearchResponse(data)
	}

	// Build query URL
	baseURL := c.buildURL(req.EntityType)
	params := url.Values{}

	// Build WHERE clause from filters
	where := query.BuildWhereClause(req.Filters, req.RawWhere)
	if where != "" {
		params.Set("where", where)
	}

	// Include fields
	if len(req.Include) > 0 {
		params.Set("include", fmt.Sprintf("[%s]", strings.Join(req.Include, ",")))
	}

	// Take (limit)
	if req.Take > 0 {
		params.Set("take", fmt.Sprintf("%d", req.Take))
	}

	// OrderBy
	if req.OrderByField != "" {
		if req.OrderByDesc {
			params.Set("orderByDesc", req.OrderByField)
		} else {
			params.Set("orderBy", req.OrderByField)
		}
	}

	fullURL := baseURL
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	data, err := c.doGet(ctx, fullURL)
	if err != nil {
		return nil, err
	}
	return c.parseSearchResponse(data)
}

func (c *httpClient) parseSearchResponse(data []byte) (*query.PaginatedResponse, error) {
	var apiResp entity.APIResponse
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, err
	}

	hasMore := apiResp.Next != ""
	return &query.PaginatedResponse{
		Items: apiResp.Items,
		Pagination: query.PaginationMeta{
			HasMore:  hasMore,
			Cursor:   apiResp.Next,
			Returned: len(apiResp.Items),
		},
	}, nil
}
