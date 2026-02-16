package testutil

import (
	"fmt"

	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/domain/query"
)

// NewSearchResponse creates a test PaginatedResponse
func NewSearchResponse(items int) *query.PaginatedResponse {
	result := &query.PaginatedResponse{
		Items: make([]map[string]any, items),
		Pagination: query.PaginationMeta{
			HasMore:  false,
			Returned: items,
		},
	}
	for i := 0; i < items; i++ {
		result.Items[i] = map[string]any{
			"Id":           float64(i + 1),
			"Name":         fmt.Sprintf("Item %d", i+1),
			"ResourceType": "UserStory",
		}
	}
	return result
}

// NewEntityMap creates a test entity map
func NewEntityMap(entityType, name string) map[string]any {
	return map[string]any{
		"Id":           float64(1),
		"Name":         name,
		"ResourceType": entityType,
	}
}

// NewComment creates a test Comment
func NewComment(id int, description string) entity.Comment {
	return entity.Comment{
		ID:          id,
		Description: description,
		CreateDate:  "2024-01-01T00:00:00",
		IsPrivate:   true,
	}
}

// NewAttachment creates a test Attachment
func NewAttachment(id int, name string) entity.Attachment {
	mime := "application/pdf"
	return entity.Attachment{
		ID:       id,
		Name:     name,
		Size:     1024,
		MimeType: &mime,
		Uri:      "/api/v1/Attachments/1/download",
	}
}
