package docs

import (
	fxctx "github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
)

func ptr(v string) *string {
	return &v
}

type resourceDef struct {
	URI         string
	Name        string
	Description string
	Content     string
}

func Resources() []fxctx.Resource {
	defs := []resourceDef{
		{"docs://getting-started", "Getting Started Guide", "Quick start guide for using the TP MCP server", GettingStarted},
		{"docs://tool-reference", "Tool Reference", "Complete reference for all MCP tools and parameters", ToolReference},
		{"docs://examples", "Usage Examples", "Comprehensive usage examples organized by scenario", Examples},
		{"docs://query-guide", "Query Guide", "Guide to WHERE clause syntax and date macros", QueryGuide},
		{"docs://authentication", "Authentication Guide", "Authentication configuration reference", Authentication},
	}

	resources := make([]fxctx.Resource, len(defs))
	for i, def := range defs {
		content := def.Content // capture for closure
		uri := def.URI
		resources[i] = fxctx.NewResource(
			mcp.Resource{
				Uri:         def.URI,
				Name:        def.Name,
				Description: ptr(def.Description),
				MimeType:    ptr("text/markdown"),
			},
			func(requestURI string) (*mcp.ReadResourceResult, error) {
				return &mcp.ReadResourceResult{
					Contents: []interface{}{
						mcp.TextResourceContents{
							Uri:      uri,
							MimeType: ptr("text/markdown"),
							Text:     content,
						},
					},
				}, nil
			},
		)
	}
	return resources
}
