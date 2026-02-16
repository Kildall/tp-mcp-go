package tools

import (
	"context"
	"encoding/json"
	"testing"

	"tp-mcp-go/internal/docs"
	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/domain/query"
	"tp-mcp-go/internal/testutil"

	fxctx "github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
)

func TestAllToolsRegistered(t *testing.T) {
	mock := &testutil.MockClient{
		GetValidEntityTypesFn: func(ctx context.Context) ([]string, error) {
			return []string{"UserStory", "Bug"}, nil
		},
	}

	tools := []struct {
		name string
		tool fxctx.Tool
	}{
		{"search", NewSearchTool(mock)},
		{"get_entity", NewGetEntityTool(mock)},
		{"create_entity", NewCreateEntityTool(mock)},
		{"update_entity", NewUpdateEntityTool(mock)},
		{"add_comment", NewAddCommentTool(mock)},
		{"list_comments", NewListCommentsTool(mock)},
		{"list_attachments", NewListAttachmentsTool(mock)},
		{"download_attachment", NewDownloadAttachmentTool(mock)},
		{"inspect_object", NewInspectObjectTool(mock)},
		{"get_documentation", NewGetDocumentationTool()},
	}

	for _, tt := range tools {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tool.GetMcpTool().Name != tt.name {
				t.Errorf("expected tool name %q, got %q", tt.name, tt.tool.GetMcpTool().Name)
			}
		})
	}
}

func TestSearchToolIntegration(t *testing.T) {
	mock := &testutil.MockClient{
		SearchEntitiesFn: func(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error) {
			return testutil.NewSearchResponse(2), nil
		},
	}

	tool := NewSearchTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type": "UserStory",
	})

	if result.IsError != nil && *result.IsError {
		t.Fatalf("expected success, got error: %v", result.Content)
	}

	// Verify response has JSON content
	if len(result.Content) == 0 {
		t.Fatal("expected non-empty content")
	}

	// Verify we can parse the JSON response
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(textContent.Text), &data); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}
	// Verify it has items
	if items, ok := data["items"].([]interface{}); !ok || len(items) == 0 {
		t.Error("expected items in response")
	}
}

func TestGetDocumentationIntegration(t *testing.T) {
	tool := NewGetDocumentationTool()
	result := tool.Callback(map[string]interface{}{
		"topic": "search",
	})

	if result.IsError != nil && *result.IsError {
		t.Fatalf("expected success, got error: %v", result.Content)
	}

	// Verify response has text content
	if len(result.Content) == 0 {
		t.Fatal("expected non-empty content")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	if len(textContent.Text) == 0 {
		t.Error("expected non-empty documentation text")
	}
}

func TestGetEntityIntegration(t *testing.T) {
	mock := &testutil.MockClient{
		GetEntityFn: func(ctx context.Context, entityType entity.Type, id int, include []string) (map[string]any, error) {
			return testutil.NewEntityMap("UserStory", "Test Story"), nil
		},
	}

	tool := NewGetEntityTool(mock)
	result := tool.Callback(map[string]interface{}{
		"type": "UserStory",
		"id":   1,
	})

	if result.IsError != nil && *result.IsError {
		t.Fatalf("expected success, got error: %v", result.Content)
	}

	// Verify response has JSON content
	if len(result.Content) == 0 {
		t.Fatal("expected non-empty content")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(textContent.Text), &data); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}
	// Verify entity has expected fields
	if name, ok := data["Name"].(string); !ok || name != "Test Story" {
		t.Error("expected Name field with value 'Test Story'")
	}
}

func TestDocumentationAllTopics(t *testing.T) {
	tool := NewGetDocumentationTool()

	topics := []string{
		"overview",
		"tools",
		"search",
		"entities",
		"comments",
		"attachments",
		"inspect",
		"authentication",
		"pagination",
		"query-syntax",
		"examples",
	}

	for _, topic := range topics {
		t.Run(topic, func(t *testing.T) {
			result := tool.Callback(map[string]interface{}{
				"topic": topic,
			})

			if result.IsError != nil && *result.IsError {
				t.Fatalf("expected success for topic %q, got error: %v", topic, result.Content)
			}

			// Verify response has non-empty content
			if len(result.Content) == 0 {
				t.Fatalf("expected non-empty content for topic %q", topic)
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatalf("expected mcp.TextContent for topic %q, got %T", topic, result.Content[0])
			}

			if len(textContent.Text) == 0 {
				t.Errorf("expected non-empty documentation text for topic %q", topic)
			}
		})
	}
}

func TestInspectListTypes(t *testing.T) {
	mock := &testutil.MockClient{
		GetValidEntityTypesFn: func(ctx context.Context) ([]string, error) {
			return []string{"UserStory", "Bug", "Task", "Feature"}, nil
		},
	}

	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action": "list_types",
	})

	if result.IsError != nil && *result.IsError {
		t.Fatalf("expected success, got error: %v", result.Content)
	}

	// Verify response has JSON content
	if len(result.Content) == 0 {
		t.Fatal("expected non-empty content")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	var types []string
	if err := json.Unmarshal([]byte(textContent.Text), &types); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}
	// Verify it returns a list of entity types
	if len(types) == 0 {
		t.Error("expected non-empty list of entity types")
	}
	if len(types) != 4 {
		t.Errorf("expected 4 types, got %d", len(types))
	}
}

func TestResourcesExist(t *testing.T) {
	resources := docs.Resources()

	if len(resources) != 5 {
		t.Errorf("expected 5 resources, got %d", len(resources))
	}

	expectedURIs := map[string]bool{
		"docs://getting-started": false,
		"docs://tool-reference":  false,
		"docs://examples":        false,
		"docs://query-guide":     false,
		"docs://authentication":  false,
	}

	for _, res := range resources {
		uri := res.GetResource().Uri
		if _, ok := expectedURIs[uri]; ok {
			expectedURIs[uri] = true
		} else {
			t.Errorf("unexpected resource URI: %s", uri)
		}
	}

	for uri, found := range expectedURIs {
		if !found {
			t.Errorf("expected resource URI %s not found", uri)
		}
	}
}
