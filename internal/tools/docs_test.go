package tools

import (
	"testing"

	"github.com/strowk/foxy-contexts/pkg/mcp"
)

func TestGetDocumentationTool(t *testing.T) {
	tool := NewGetDocumentationTool()
	if tool == nil {
		t.Fatal("expected tool to be created")
	}

	t.Run("valid topic returns content", func(t *testing.T) {
		result := tool.Callback(map[string]interface{}{
			"topic": "search",
		})

		if result == nil {
			t.Fatal("expected result to not be nil")
		}
		if result.IsError != nil && *result.IsError {
			t.Fatal("expected no error")
		}
		if len(result.Content) != 1 {
			t.Fatalf("expected 1 content item, got %d", len(result.Content))
		}

		// Extract text content
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
		}

		if len(textContent.Text) == 0 {
			t.Error("expected non-empty content")
		}
		if textContent.Text[:14] != "# Search Guide" {
			t.Error("expected content to start with '# Search Guide'")
		}
	})

	t.Run("overview topic returns getting started with topic list", func(t *testing.T) {
		result := tool.Callback(map[string]interface{}{
			"topic": "overview",
		})

		if result == nil {
			t.Fatal("expected result to not be nil")
		}
		if result.IsError != nil && *result.IsError {
			t.Fatal("expected no error")
		}
		if len(result.Content) != 1 {
			t.Fatalf("expected 1 content item, got %d", len(result.Content))
		}

		// Extract text content
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
		}

		text := textContent.Text
		if len(text) == 0 {
			t.Error("expected non-empty content")
		}
		// Check for getting started content
		if text[:25] != "# TP MCP Server - Getting" {
			t.Error("expected content to contain getting started guide")
		}
		// Should also contain available topics list
		hasTopicsList := false
		for i := 0; i < len(text)-20; i++ {
			if text[i:i+20] == "## Available Topics\n" {
				hasTopicsList = true
				break
			}
		}
		if !hasTopicsList {
			t.Error("expected content to contain available topics list")
		}
	})

	t.Run("invalid topic returns overview with topic list", func(t *testing.T) {
		result := tool.Callback(map[string]interface{}{
			"topic": "nonexistent-topic",
		})

		if result == nil {
			t.Fatal("expected result to not be nil")
		}
		if result.IsError != nil && *result.IsError {
			t.Fatal("expected no error")
		}
		if len(result.Content) != 1 {
			t.Fatalf("expected 1 content item, got %d", len(result.Content))
		}

		// Extract text content
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("expected mcp.TextContent, got %T", result.Content[0])
		}

		text := textContent.Text
		if len(text) == 0 {
			t.Error("expected non-empty content")
		}
		// Should return getting started with topics list
		if text[:25] != "# TP MCP Server - Getting" {
			t.Error("expected content to contain getting started guide")
		}
	})

	t.Run("all 11 topic keys return non-empty content", func(t *testing.T) {
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

				if result == nil {
					t.Fatalf("result should not be nil for topic: %s", topic)
				}
				if result.IsError != nil && *result.IsError {
					t.Fatalf("should not have error for topic: %s", topic)
				}
				if len(result.Content) != 1 {
					t.Fatalf("should have exactly 1 content item for topic: %s, got %d", topic, len(result.Content))
				}

				// Extract text content
				textContent, ok := result.Content[0].(mcp.TextContent)
				if !ok {
					t.Fatalf("expected mcp.TextContent for topic: %s, got %T", topic, result.Content[0])
				}

				if len(textContent.Text) == 0 {
					t.Errorf("content should not be empty for topic: %s", topic)
				}
			})
		}
	})
}
