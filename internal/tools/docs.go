package tools

import (
	"fmt"
	"strings"

	"tp-mcp-go/internal/docs"

	fxctx "github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
)

// NewGetDocumentationTool creates the get_documentation tool
// NOTE: NO client dependency — uses embedded docs only
func NewGetDocumentationTool() fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name:        "get_documentation",
			Description: ptr("Access embedded documentation for the TP MCP server"),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"topic": {
						"type":        "string",
						"description": "Documentation topic to retrieve",
						"enum":        topicKeysList(),
					},
				},
				Required: []string{"topic"},
			},
		},
		func(args map[string]interface{}) *mcp.CallToolResult {
			topic := getStringArg(args, "topic")

			// For "overview" or invalid topic, return getting started with topic list
			if topic == "overview" || topic == "" {
				return textResult(buildOverviewWithTopics())
			}

			doc, err := docs.GetTopic(topic)
			if err != nil {
				// Return overview with topic list for invalid topics
				return textResult(buildOverviewWithTopics())
			}

			return textResult(doc.Content)
		},
	)
}

func topicKeysList() []interface{} {
	topics := docs.ListTopics()
	keys := make([]interface{}, len(topics))
	for i, t := range topics {
		keys[i] = t.Key
	}
	return keys
}

func buildOverviewWithTopics() string {
	var sb strings.Builder
	sb.WriteString(docs.GettingStarted)
	sb.WriteString("\n\n## Available Topics\n\n")
	for _, t := range docs.ListTopics() {
		sb.WriteString(fmt.Sprintf("- **%s**: %s\n", t.Key, t.Description))
	}
	return sb.String()
}
