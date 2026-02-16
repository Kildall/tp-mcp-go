package client

import (
	"context"
	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/domain/query"
)

type Client interface {
	// Search
	SearchEntities(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error)

	// Entity CRUD
	GetEntity(ctx context.Context, entityType entity.Type, id int, include []string) (map[string]any, error)
	CreateEntity(ctx context.Context, entityType entity.Type, data map[string]any) (map[string]any, error)
	UpdateEntity(ctx context.Context, entityType entity.Type, id int, data map[string]any) (map[string]any, error)

	// Comments — note: ListComments has an include parameter
	CreateComment(ctx context.Context, entityID int, description string) (*entity.Comment, error)
	ListComments(ctx context.Context, entityID int, take int, include []string) ([]entity.Comment, error)

	// Attachments — note: GetAttachmentMetadata for size-guard support
	ListAttachments(ctx context.Context, entityID int, take int) ([]entity.Attachment, error)
	GetAttachmentMetadata(ctx context.Context, attachmentID int) (*entity.Attachment, error)
	DownloadAttachment(ctx context.Context, uri string) ([]byte, string, error) // bytes, mimeType, error

	// Metadata
	FetchMetadata(ctx context.Context) (any, error)
	GetValidEntityTypes(ctx context.Context) ([]string, error)
	InitializeCache(ctx context.Context) error
}
