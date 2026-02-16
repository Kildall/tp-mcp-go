package testutil

import (
	"context"
	"tp-mcp-go/internal/client"
	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/domain/query"
)

// MockClient satisfies client.Client with configurable function fields
type MockClient struct {
	SearchEntitiesFn        func(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error)
	GetEntityFn             func(ctx context.Context, entityType entity.Type, id int, include []string) (map[string]any, error)
	CreateEntityFn          func(ctx context.Context, entityType entity.Type, data map[string]any) (map[string]any, error)
	UpdateEntityFn          func(ctx context.Context, entityType entity.Type, id int, data map[string]any) (map[string]any, error)
	CreateCommentFn         func(ctx context.Context, entityID int, description string) (*entity.Comment, error)
	ListCommentsFn          func(ctx context.Context, entityID int, take int, include []string) ([]entity.Comment, error)
	ListAttachmentsFn       func(ctx context.Context, entityID int, take int) ([]entity.Attachment, error)
	GetAttachmentMetadataFn func(ctx context.Context, attachmentID int) (*entity.Attachment, error)
	DownloadAttachmentFn    func(ctx context.Context, uri string) ([]byte, string, error)
	FetchMetadataFn         func(ctx context.Context) (any, error)
	GetValidEntityTypesFn   func(ctx context.Context) ([]string, error)
	InitializeCacheFn       func(ctx context.Context) error
}

// Compile-time interface check
var _ client.Client = (*MockClient)(nil)

func (m *MockClient) SearchEntities(ctx context.Context, req query.SearchRequest) (*query.PaginatedResponse, error) {
	if m.SearchEntitiesFn != nil {
		return m.SearchEntitiesFn(ctx, req)
	}
	return &query.PaginatedResponse{}, nil
}

func (m *MockClient) GetEntity(ctx context.Context, entityType entity.Type, id int, include []string) (map[string]any, error) {
	if m.GetEntityFn != nil {
		return m.GetEntityFn(ctx, entityType, id, include)
	}
	return map[string]any{}, nil
}

func (m *MockClient) CreateEntity(ctx context.Context, entityType entity.Type, data map[string]any) (map[string]any, error) {
	if m.CreateEntityFn != nil {
		return m.CreateEntityFn(ctx, entityType, data)
	}
	return map[string]any{}, nil
}

func (m *MockClient) UpdateEntity(ctx context.Context, entityType entity.Type, id int, data map[string]any) (map[string]any, error) {
	if m.UpdateEntityFn != nil {
		return m.UpdateEntityFn(ctx, entityType, id, data)
	}
	return map[string]any{}, nil
}

func (m *MockClient) CreateComment(ctx context.Context, entityID int, description string) (*entity.Comment, error) {
	if m.CreateCommentFn != nil {
		return m.CreateCommentFn(ctx, entityID, description)
	}
	return &entity.Comment{}, nil
}

func (m *MockClient) ListComments(ctx context.Context, entityID int, take int, include []string) ([]entity.Comment, error) {
	if m.ListCommentsFn != nil {
		return m.ListCommentsFn(ctx, entityID, take, include)
	}
	return []entity.Comment{}, nil
}

func (m *MockClient) ListAttachments(ctx context.Context, entityID int, take int) ([]entity.Attachment, error) {
	if m.ListAttachmentsFn != nil {
		return m.ListAttachmentsFn(ctx, entityID, take)
	}
	return []entity.Attachment{}, nil
}

func (m *MockClient) GetAttachmentMetadata(ctx context.Context, attachmentID int) (*entity.Attachment, error) {
	if m.GetAttachmentMetadataFn != nil {
		return m.GetAttachmentMetadataFn(ctx, attachmentID)
	}
	return &entity.Attachment{}, nil
}

func (m *MockClient) DownloadAttachment(ctx context.Context, uri string) ([]byte, string, error) {
	if m.DownloadAttachmentFn != nil {
		return m.DownloadAttachmentFn(ctx, uri)
	}
	return []byte{}, "", nil
}

func (m *MockClient) FetchMetadata(ctx context.Context) (any, error) {
	if m.FetchMetadataFn != nil {
		return m.FetchMetadataFn(ctx)
	}
	return nil, nil
}

func (m *MockClient) GetValidEntityTypes(ctx context.Context) ([]string, error) {
	if m.GetValidEntityTypesFn != nil {
		return m.GetValidEntityTypesFn(ctx)
	}
	return []string{}, nil
}

func (m *MockClient) InitializeCache(ctx context.Context) error {
	if m.InitializeCacheFn != nil {
		return m.InitializeCacheFn(ctx)
	}
	return nil
}
