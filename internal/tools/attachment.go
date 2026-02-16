package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"tp-mcp-go/internal/client"

	fxctx "github.com/strowk/foxy-contexts/pkg/fxctx"
	"github.com/strowk/foxy-contexts/pkg/mcp"
)

// NewListAttachmentsTool creates a tool to list attachments for an entity
func NewListAttachmentsTool(c client.Client) fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name: "list_attachments",
			Description: ptr("List attachments for a Target Process entity. " +
				"Returns attachment metadata including ID, name, size, and MIME type."),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"entityId": {
						"type":        "integer",
						"description": "Entity ID to list attachments for",
					},
					"take": {
						"type":        "integer",
						"description": "Number of attachments to return (default: 100)",
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

			// Parse take with default
			take := 100 // default
			if t, err := getIntArg(args, "take"); err == nil {
				take = t
			}

			// Call client
			attachments, err := c.ListAttachments(context.Background(), entityID, take)
			if err != nil {
				return errorResult(err)
			}

			return jsonResult(attachments)
		},
	)
}

// NewDownloadAttachmentTool creates a tool to download an attachment
func NewDownloadAttachmentTool(c client.Client) fxctx.Tool {
	return fxctx.NewTool(
		&mcp.Tool{
			Name: "download_attachment",
			Description: ptr("Download a Target Process attachment by ID. " +
				"Returns the attachment content with appropriate encoding based on MIME type. " +
				"Images are returned as base64-encoded image content, text files as plain text, " +
				"and other types as base64 with a note. " +
				"Maximum file size: 50MB."),
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]map[string]interface{}{
					"attachmentId": {
						"type":        "integer",
						"description": "Attachment ID to download",
					},
				},
				Required: []string{"attachmentId"},
			},
		},
		func(args map[string]interface{}) *mcp.CallToolResult {
			// Parse attachmentId
			attachmentID, err := getIntArg(args, "attachmentId")
			if err != nil {
				return errorResult(err)
			}

			// Get attachment metadata
			attachment, err := c.GetAttachmentMetadata(context.Background(), attachmentID)
			if err != nil {
				return errorResult(err)
			}

			// Check size limit (50MB)
			const maxSize = 50 * 1024 * 1024
			if attachment.Size > maxSize {
				return errorResult(fmt.Errorf("attachment size (%d bytes) exceeds maximum allowed size (50MB)", attachment.Size))
			}

			// Download attachment
			downloadUri := fmt.Sprintf("/Attachment.aspx?AttachmentID=%d", attachmentID)
		data, mimeType, err := c.DownloadAttachment(context.Background(), downloadUri)
			if err != nil {
				return errorResult(err)
			}

			// Route based on MIME type
			if strings.HasPrefix(mimeType, "image/") {
				// Return as ImageContent
				base64Data := base64.StdEncoding.EncodeToString(data)
				return &mcp.CallToolResult{
					Content: []interface{}{
						mcp.ImageContent{
							Type:     "image",
							Data:     base64Data,
							MimeType: mimeType,
						},
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("Filename: %s\nSize: %d bytes\nMIME Type: %s", attachment.Name, len(data), mimeType),
						},
					},
				}
			} else if strings.HasPrefix(mimeType, "text/") {
				// Return as plain text
				return &mcp.CallToolResult{
					Content: []interface{}{
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("Filename: %s\nSize: %d bytes\nMIME Type: %s\n\n%s", attachment.Name, len(data), mimeType, string(data)),
						},
					},
				}
			} else {
				// Return as base64 with note
				base64Data := base64.StdEncoding.EncodeToString(data)
				return &mcp.CallToolResult{
					Content: []interface{}{
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("Filename: %s\nSize: %d bytes\nMIME Type: %s\n\nBase64-encoded content:\n%s", attachment.Name, len(data), mimeType, base64Data),
						},
					},
				}
			}
		},
	)
}
