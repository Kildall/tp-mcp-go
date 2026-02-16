package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/testutil"

	"github.com/strowk/foxy-contexts/pkg/mcp"
	"github.com/stretchr/testify/assert"
)

func TestListAttachments(t *testing.T) {
	expectedAttachments := []entity.Attachment{
		{ID: 1, Name: "file1.png", Size: 1024},
		{ID: 2, Name: "file2.pdf", Size: 2048},
	}

	mockClient := &testutil.MockClient{
		ListAttachmentsFn: func(ctx context.Context, entityID int, take int) ([]entity.Attachment, error) {
			assert.Equal(t, 300, entityID)
			assert.Equal(t, 100, take)
			return expectedAttachments, nil
		},
	}

	tool := NewListAttachmentsTool(mockClient)
	result := tool.Callback(map[string]interface{}{
		"entityId": float64(300),
	})

	assert.NotNil(t, result)
	assert.Nil(t, result.IsError)
}

func TestDownloadAttachment_Image(t *testing.T) {
	imageData := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header
	attachment := &entity.Attachment{
		ID:   1,
		Name: "test.png",
		Size: len(imageData),
		Uri:  "http://example.com/test.png",
	}

	mockClient := &testutil.MockClient{
		GetAttachmentMetadataFn: func(ctx context.Context, attachmentID int) (*entity.Attachment, error) {
			assert.Equal(t, 1, attachmentID)
			return attachment, nil
		},
		DownloadAttachmentFn: func(ctx context.Context, uri string) ([]byte, string, error) {
			assert.Equal(t, "/Attachment.aspx?AttachmentID=1", uri)
			return imageData, "image/png", nil
		},
	}

	tool := NewDownloadAttachmentTool(mockClient)
	result := tool.Callback(map[string]interface{}{
		"attachmentId": float64(1),
	})

	assert.NotNil(t, result)
	assert.Nil(t, result.IsError)
	assert.Len(t, result.Content, 2)

	// Verify ImageContent is returned
	imageContent, ok := result.Content[0].(mcp.ImageContent)
	if !ok {
		t.Fatalf("Expected first content to be mcp.ImageContent, got %T", result.Content[0])
	}
	assert.Equal(t, "image", imageContent.Type)
	assert.Equal(t, base64.StdEncoding.EncodeToString(imageData), imageContent.Data)
	assert.Equal(t, "image/png", imageContent.MimeType)
}

func TestDownloadAttachment_Text(t *testing.T) {
	textData := []byte("Hello, World!")
	attachment := &entity.Attachment{
		ID:   2,
		Name: "test.txt",
		Size: len(textData),
		Uri:  "http://example.com/test.txt",
	}

	mockClient := &testutil.MockClient{
		GetAttachmentMetadataFn: func(ctx context.Context, attachmentID int) (*entity.Attachment, error) {
			assert.Equal(t, 2, attachmentID)
			return attachment, nil
		},
		DownloadAttachmentFn: func(ctx context.Context, uri string) ([]byte, string, error) {
			assert.Equal(t, "/Attachment.aspx?AttachmentID=2", uri)
			return textData, "text/plain", nil
		},
	}

	tool := NewDownloadAttachmentTool(mockClient)
	result := tool.Callback(map[string]interface{}{
		"attachmentId": float64(2),
	})

	assert.NotNil(t, result)
	assert.Nil(t, result.IsError)
	assert.Len(t, result.Content, 1)

	// Verify text content
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected content to be mcp.TextContent, got %T", result.Content[0])
	}
	assert.Equal(t, "text", textContent.Type)
	assert.Contains(t, textContent.Text, "Hello, World!")
}

func TestDownloadAttachment_SizeGuard(t *testing.T) {
	// Create attachment larger than 50MB
	largeSize := 51 * 1024 * 1024
	attachment := &entity.Attachment{
		ID:   3,
		Name: "large.bin",
		Size: largeSize,
		Uri:  "http://example.com/large.bin",
	}

	mockClient := &testutil.MockClient{
		GetAttachmentMetadataFn: func(ctx context.Context, attachmentID int) (*entity.Attachment, error) {
			assert.Equal(t, 3, attachmentID)
			return attachment, nil
		},
	}

	tool := NewDownloadAttachmentTool(mockClient)
	result := tool.Callback(map[string]interface{}{
		"attachmentId": float64(3),
	})

	assert.NotNil(t, result)
	assert.NotNil(t, result.IsError)
	assert.True(t, *result.IsError)

	// Verify error message
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected content to be mcp.TextContent, got %T", result.Content[0])
	}
	assert.Contains(t, textContent.Text, fmt.Sprintf("Error: attachment size (%d bytes) exceeds maximum allowed size (50MB)", largeSize))
}
