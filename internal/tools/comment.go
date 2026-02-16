package tools

import (
	"context"
	"fmt"

	"tp-mcp-go/internal/client"

	fxctx "github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
)

// NewAddCommentTool creates a tool to add a comment to an entity
func NewAddCommentTool(c client.Client) fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name: "add_comment",
			Description: ptr("Add a comment to a Target Process entity. " +
				"Supports HTML formatting in the description."),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"entityId": {
						"type":        "integer",
						"description": "Entity ID to add the comment to",
					},
					"description": {
						"type":        "string",
						"description": "Comment text (supports HTML formatting)",
					},
				},
				Required: []string{"entityId", "description"},
			},
		},
		func(args map[string]interface{}) *mcp.CallToolResult {
			// Parse entityId
			entityID, err := getIntArg(args, "entityId")
			if err != nil {
				return errorResult(err)
			}

			// Parse description
			description := getStringArg(args, "description")
			if description == "" {
				return errorResult(fmt.Errorf("description parameter is required"))
			}

			// Call client
			comment, err := c.CreateComment(context.Background(), entityID, description)
			if err != nil {
				return errorResult(err)
			}

			return jsonResult(comment)
		},
	)
}

// NewListCommentsTool creates a tool to list comments for an entity
func NewListCommentsTool(c client.Client) fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name: "list_comments",
			Description: ptr("List comments for a Target Process entity. " +
				"Returns paginated comments with optional field inclusion."),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"entityId": {
						"type":        "integer",
						"description": "Entity ID to list comments for",
					},
					"take": {
						"type":        "integer",
						"description": "Number of comments to return (default: 25, min: 1, max: 100)",
						"minimum":     1,
						"maximum":     100,
					},
					"include": {
						"type":        "array",
						"description": "Fields to include in response (default: [Description,CreateDate,Owner])",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
				Required: []string{"entityId"},
			},
		},
		func(args map[string]interface{}) *mcp.CallToolResult {
			// Parse entityId
			entityID, err := getIntArg(args, "entityId")
			if err != nil {
				return errorResult(err)
			}

			// Parse take with default and clamping
			take := 25 // default
			if t, err := getIntArg(args, "take"); err == nil {
				take = t
			}
			// Clamp to [1, 100]
			if take < 1 {
				take = 1
			}
			if take > 100 {
				take = 100
			}

			// Parse include with default
			include := getStringSliceArg(args, "include")
			if include == nil {
				include = []string{"Description", "CreateDate", "Owner"}
			}

			// Call client
			comments, err := c.ListComments(context.Background(), entityID, take, include)
			if err != nil {
				return errorResult(err)
			}

			return jsonResult(comments)
		},
	)
}
