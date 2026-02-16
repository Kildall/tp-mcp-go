package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"tp-mcp-go/internal/domain/query"
	"tp-mcp-go/internal/testutil"

	"github.com/strowk/foxy-contexts/pkg/mcp"
)

func TestSearchToolReturnsResults(t *testing.T) {
	mock := &testutil.MockClient{
		SearchEntitiesFn: func(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error) {
			return testutil.NewSearchResponse(3), nil
		},
	}

	tool := NewSearchTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type": "UserStory",
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	// Verify it returns JSON with 3 items
	if len(result.Content) == 0 {
		t.Fatal("expected content in result")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	text := textContent.Text

	var response query.PaginatedResponse
	if err := json.Unmarshal([]byte(text), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(response.Items))
	}

	if response.Pagination.Returned != 3 {
		t.Errorf("expected pagination.returned = 3, got %d", response.Pagination.Returned)
	}
}

func TestSearchToolHandlesError(t *testing.T) {
	mock := &testutil.MockClient{
		SearchEntitiesFn: func(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	tool := NewSearchTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type": "UserStory",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error")
	}

	// Verify error message is in content
	if len(result.Content) == 0 {
		t.Fatal("expected content in error result")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	if textContent.Text == "" {
		t.Error("expected error message in text")
	}
}

func TestSearchToolClampsTake(t *testing.T) {
	var capturedReq query.SearchRequest
	mock := &testutil.MockClient{
		SearchEntitiesFn: func(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error) {
			capturedReq = req
			return testutil.NewSearchResponse(0), nil
		},
	}

	tool := NewSearchTool(mock)

	// Test take too high
	tool.Callback(map[string]interface{}{
		"type": "UserStory",
		"take": float64(5000),
	})
	if capturedReq.Take != 1000 {
		t.Errorf("expected take clamped to 1000, got %d", capturedReq.Take)
	}

	// Test take too low
	tool.Callback(map[string]interface{}{
		"type": "UserStory",
		"take": float64(0),
	})
	if capturedReq.Take != 1 {
		t.Errorf("expected take clamped to 1, got %d", capturedReq.Take)
	}

	// Test take within range
	tool.Callback(map[string]interface{}{
		"type": "UserStory",
		"take": float64(50),
	})
	if capturedReq.Take != 50 {
		t.Errorf("expected take = 50, got %d", capturedReq.Take)
	}

	// Test default take when not provided
	tool.Callback(map[string]interface{}{
		"type": "UserStory",
	})
	if capturedReq.Take != 100 {
		t.Errorf("expected default take = 100, got %d", capturedReq.Take)
	}
}

func TestSearchToolPaginationMetadata(t *testing.T) {
	mock := &testutil.MockClient{
		SearchEntitiesFn: func(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error) {
			resp := testutil.NewSearchResponse(2)
			resp.Pagination.HasMore = true
			resp.Pagination.Cursor = "https://example.com/api/v1/UserStorys?next=abc"
			return resp, nil
		},
	}

	tool := NewSearchTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type": "UserStory",
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	// Extract and verify pagination data
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	text := textContent.Text

	var response query.PaginatedResponse
	if err := json.Unmarshal([]byte(text), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !response.Pagination.HasMore {
		t.Error("expected hasMore = true")
	}

	if response.Pagination.Cursor != "https://example.com/api/v1/UserStorys?next=abc" {
		t.Errorf("expected cursor to be set, got %q", response.Pagination.Cursor)
	}

	if response.Pagination.Returned != 2 {
		t.Errorf("expected returned = 2, got %d", response.Pagination.Returned)
	}
}

func TestSearchToolInvalidType(t *testing.T) {
	mock := &testutil.MockClient{}
	tool := NewSearchTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type": "InvalidType",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for invalid entity type")
	}

	// Verify error message mentions invalid type
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	if textContent.Text == "" {
		t.Error("expected error message in text")
	}
}

func TestSearchToolWithFilters(t *testing.T) {
	var capturedReq query.SearchRequest
	mock := &testutil.MockClient{
		SearchEntitiesFn: func(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error) {
			capturedReq = req
			return testutil.NewSearchResponse(1), nil
		},
	}

	tool := NewSearchTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type":         "UserStory",
		"status":       "Open",
		"assignedUser": "user@example.com",
		"priority":     "High",
		"dateFrom":     "2024-01-01",
		"dateTo":       "2024-12-31",
		"dateField":    "CreateDate",
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	// Verify filters were captured correctly
	if capturedReq.Filters.Status != "Open" {
		t.Errorf("expected status = Open, got %q", capturedReq.Filters.Status)
	}

	if capturedReq.Filters.AssignedUser != "user@example.com" {
		t.Errorf("expected assignedUser = user@example.com, got %q", capturedReq.Filters.AssignedUser)
	}

	if capturedReq.Filters.Priority != "High" {
		t.Errorf("expected priority = High, got %q", capturedReq.Filters.Priority)
	}

	if capturedReq.Filters.DateFrom != "2024-01-01" {
		t.Errorf("expected dateFrom = 2024-01-01, got %q", capturedReq.Filters.DateFrom)
	}

	if capturedReq.Filters.DateTo != "2024-12-31" {
		t.Errorf("expected dateTo = 2024-12-31, got %q", capturedReq.Filters.DateTo)
	}

	if capturedReq.Filters.DateField != "CreateDate" {
		t.Errorf("expected dateField = CreateDate, got %q", capturedReq.Filters.DateField)
	}
}

func TestSearchToolWithCursor(t *testing.T) {
	var capturedReq query.SearchRequest
	mock := &testutil.MockClient{
		SearchEntitiesFn: func(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error) {
			capturedReq = req
			return testutil.NewSearchResponse(2), nil
		},
	}

	tool := NewSearchTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type":   "UserStory",
		"cursor": "https://example.com/api/v1/UserStorys?next=abc",
		"status": "Open", // Should be ignored when cursor is provided
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	// Verify cursor was set
	if capturedReq.Cursor != "https://example.com/api/v1/UserStorys?next=abc" {
		t.Errorf("expected cursor to be set, got %q", capturedReq.Cursor)
	}

	// Verify other filters are ignored when cursor is provided
	if capturedReq.Filters.Status != "" {
		t.Error("expected filters to be empty when cursor is provided")
	}
}
