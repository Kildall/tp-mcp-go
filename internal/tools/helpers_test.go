package tools

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/strowk/foxy-contexts/pkg/mcp"
)

func TestGetStringArg(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		key      string
		expected string
	}{
		{
			name:     "returns value when present",
			args:     map[string]any{"name": "test"},
			key:      "name",
			expected: "test",
		},
		{
			name:     "returns empty when missing",
			args:     map[string]any{},
			key:      "name",
			expected: "",
		},
		{
			name:     "returns empty when wrong type",
			args:     map[string]any{"name": 123},
			key:      "name",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringArg(tt.args, tt.key)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetIntArg(t *testing.T) {
	tests := []struct {
		name      string
		args      map[string]any
		key       string
		expected  int
		expectErr bool
	}{
		{
			name:     "handles float64",
			args:     map[string]any{"count": float64(42)},
			key:      "count",
			expected: 42,
		},
		{
			name:     "handles int",
			args:     map[string]any{"count": 42},
			key:      "count",
			expected: 42,
		},
		{
			name:     "handles json.Number",
			args:     map[string]any{"count": json.Number("42")},
			key:      "count",
			expected: 42,
		},
		{
			name:      "returns error when missing",
			args:      map[string]any{},
			key:       "count",
			expectErr: true,
		},
		{
			name:      "returns error for non-number type",
			args:      map[string]any{"count": "not a number"},
			key:       "count",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getIntArg(tt.args, tt.key)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %d, got %d", tt.expected, result)
				}
			}
		})
	}
}

func TestGetStringSliceArg(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		key      string
		expected []string
	}{
		{
			name:     "handles array of interfaces",
			args:     map[string]any{"tags": []interface{}{"tag1", "tag2"}},
			key:      "tags",
			expected: []string{"tag1", "tag2"},
		},
		{
			name:     "returns nil when missing",
			args:     map[string]any{},
			key:      "tags",
			expected: nil,
		},
		{
			name:     "filters non-string items",
			args:     map[string]any{"tags": []interface{}{"tag1", 123, "tag2"}},
			key:      "tags",
			expected: []string{"tag1", "tag2"},
		},
		{
			name:     "returns nil when wrong type",
			args:     map[string]any{"tags": "not an array"},
			key:      "tags",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringSliceArg(tt.args, tt.key)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("at index %d: expected %q, got %q", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestErrorResult(t *testing.T) {
	err := errors.New("test error")
	result := errorResult(err)

	if result.IsError == nil || !*result.IsError {
		t.Error("expected IsError to be true")
	}

	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}

	content, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected TextContent")
	}

	if content.Type != "text" {
		t.Errorf("expected type 'text', got %q", content.Type)
	}

	if content.Text != "Error: test error" {
		t.Errorf("expected 'Error: test error', got %q", content.Text)
	}
}

func TestJsonResult(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"id":   42,
	}
	result := jsonResult(data)

	if result.IsError != nil && *result.IsError {
		t.Error("expected IsError to be false or nil")
	}

	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}

	content, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected TextContent")
	}

	if content.Type != "text" {
		t.Errorf("expected type 'text', got %q", content.Type)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(content.Text), &parsed); err != nil {
		t.Errorf("result is not valid JSON: %v", err)
	}

	if parsed["name"] != "test" || parsed["id"] != float64(42) {
		t.Errorf("JSON content doesn't match input data: %v", parsed)
	}
}

func TestTextResult(t *testing.T) {
	text := "test message"
	result := textResult(text)

	if result.IsError != nil && *result.IsError {
		t.Error("expected IsError to be false or nil")
	}

	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}

	content, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected TextContent")
	}

	if content.Type != "text" {
		t.Errorf("expected type 'text', got %q", content.Type)
	}

	if content.Text != text {
		t.Errorf("expected %q, got %q", text, content.Text)
	}
}
