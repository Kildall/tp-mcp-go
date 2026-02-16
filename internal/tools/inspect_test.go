package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"tp-mcp-go/internal/testutil"

	"github.com/strowk/foxy-contexts/pkg/mcp"
)

// mockMetadata returns a sample metadata structure
func mockMetadata() map[string]any {
	return map[string]any{
		"UserStory": map[string]any{
			"Name": map[string]any{
				"type":     "string",
				"required": true,
			},
			"Description": map[string]any{
				"type":     "string",
				"required": false,
			},
			"Priority": map[string]any{
				"type":     "object",
				"required": false,
			},
		},
		"Bug": map[string]any{
			"Name": map[string]any{
				"type":     "string",
				"required": true,
			},
			"Severity": map[string]any{
				"type":     "string",
				"required": false,
			},
		},
	}
}

func TestInspectObjectListTypes(t *testing.T) {
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
		t.Fatal("expected success, got error")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	var types []string
	if err := json.Unmarshal([]byte(textContent.Text), &types); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(types) != 4 {
		t.Errorf("expected 4 types, got %d", len(types))
	}

	if types[0] != "UserStory" {
		t.Errorf("expected first type to be UserStory, got %s", types[0])
	}
}

func TestInspectObjectGetProperties(t *testing.T) {
	mock := &testutil.MockClient{
		FetchMetadataFn: func(ctx context.Context) (any, error) {
			return mockMetadata(), nil
		},
	}

	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action":     "get_properties",
		"entityType": "UserStory",
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	var props map[string]any
	if err := json.Unmarshal([]byte(textContent.Text), &props); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if _, ok := props["Name"]; !ok {
		t.Error("expected Name property to be present")
	}

	if _, ok := props["Description"]; !ok {
		t.Error("expected Description property to be present")
	}

	if _, ok := props["Priority"]; !ok {
		t.Error("expected Priority property to be present")
	}
}

func TestInspectObjectGetPropertiesMissingEntityType(t *testing.T) {
	mock := &testutil.MockClient{}
	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action": "get_properties",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for missing entityType")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	if textContent.Text == "" {
		t.Error("expected error message in text")
	}
}

func TestInspectObjectGetPropertiesInvalidEntityType(t *testing.T) {
	mock := &testutil.MockClient{
		FetchMetadataFn: func(ctx context.Context) (any, error) {
			return mockMetadata(), nil
		},
	}

	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action":     "get_properties",
		"entityType": "InvalidType",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for invalid entity type")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	if textContent.Text == "" {
		t.Error("expected error message in text")
	}
}

func TestInspectObjectGetPropertyDetails(t *testing.T) {
	mock := &testutil.MockClient{
		FetchMetadataFn: func(ctx context.Context) (any, error) {
			return mockMetadata(), nil
		},
	}

	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action":     "get_property_details",
		"entityType": "UserStory",
		"property":   "Name",
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	var propDetails map[string]any
	if err := json.Unmarshal([]byte(textContent.Text), &propDetails); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if propDetails["type"] != "string" {
		t.Errorf("expected type to be string, got %v", propDetails["type"])
	}

	if propDetails["required"] != true {
		t.Errorf("expected required to be true, got %v", propDetails["required"])
	}
}

func TestInspectObjectGetPropertyDetailsMissingProperty(t *testing.T) {
	mock := &testutil.MockClient{
		FetchMetadataFn: func(ctx context.Context) (any, error) {
			return mockMetadata(), nil
		},
	}

	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action":     "get_property_details",
		"entityType": "UserStory",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for missing property")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	if textContent.Text == "" {
		t.Error("expected error message in text")
	}
}

func TestInspectObjectGetPropertyDetailsInvalidProperty(t *testing.T) {
	mock := &testutil.MockClient{
		FetchMetadataFn: func(ctx context.Context) (any, error) {
			return mockMetadata(), nil
		},
	}

	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action":     "get_property_details",
		"entityType": "UserStory",
		"property":   "InvalidProperty",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for invalid property")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	if textContent.Text == "" {
		t.Error("expected error message in text")
	}
}

func TestInspectObjectDiscoverAPIStructure(t *testing.T) {
	mock := &testutil.MockClient{
		FetchMetadataFn: func(ctx context.Context) (any, error) {
			return mockMetadata(), nil
		},
	}

	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action": "discover_api_structure",
	})

	if result.IsError != nil && *result.IsError {
		t.Fatal("expected success, got error")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
	}

	var metadata map[string]any
	if err := json.Unmarshal([]byte(textContent.Text), &metadata); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if _, ok := metadata["UserStory"]; !ok {
		t.Error("expected UserStory in metadata")
	}

	if _, ok := metadata["Bug"]; !ok {
		t.Error("expected Bug in metadata")
	}
}

func TestInspectObjectListTypesError(t *testing.T) {
	mock := &testutil.MockClient{
		GetValidEntityTypesFn: func(ctx context.Context) ([]string, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action": "list_types",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error")
	}
}

func TestInspectObjectGetPropertiesError(t *testing.T) {
	mock := &testutil.MockClient{
		FetchMetadataFn: func(ctx context.Context) (any, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action":     "get_properties",
		"entityType": "UserStory",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error")
	}
}

func TestInspectObjectInvalidAction(t *testing.T) {
	mock := &testutil.MockClient{}
	tool := NewInspectObjectTool(mock)
	result := tool.Callback(map[string]interface{}{
		"action": "invalid_action",
	})

	if result.IsError == nil || !*result.IsError {
		t.Fatal("expected error for invalid action")
	}
}
