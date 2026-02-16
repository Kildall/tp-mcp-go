package tools

import (
	"context"
	"testing"

	"tp-mcp-go/internal/domain/entity"
	"tp-mcp-go/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestAddComment(t *testing.T) {
	expectedComment := &entity.Comment{
		ID:          123,
		Description: "Test comment",
	}

	mockClient := &testutil.MockClient{
		CreateCommentFn: func(ctx context.Context, entityID int, description string) (*entity.Comment, error) {
			assert.Equal(t, 456, entityID)
			assert.Equal(t, "Test comment", description)
			return expectedComment, nil
		},
	}

	tool := NewAddCommentTool(mockClient)
	result := tool.Callback(map[string]interface{}{
		"entityId":    float64(456),
		"description": "Test comment",
	})

	assert.NotNil(t, result)
	assert.Nil(t, result.IsError)
}

func TestListComments_Defaults(t *testing.T) {
	expectedComments := []entity.Comment{
		{ID: 1, Description: "Comment 1"},
		{ID: 2, Description: "Comment 2"},
	}

	mockClient := &testutil.MockClient{
		ListCommentsFn: func(ctx context.Context, entityID int, take int, include []string) ([]entity.Comment, error) {
			assert.Equal(t, 100, entityID)
			assert.Equal(t, 25, take)
			assert.Equal(t, []string{"Description", "CreateDate", "Owner"}, include)
			return expectedComments, nil
		},
	}

	tool := NewListCommentsTool(mockClient)
	result := tool.Callback(map[string]interface{}{
		"entityId": float64(100),
	})

	assert.NotNil(t, result)
	assert.Nil(t, result.IsError)
}

func TestListComments_CustomParams(t *testing.T) {
	expectedComments := []entity.Comment{
		{ID: 1, Description: "Comment 1"},
	}

	mockClient := &testutil.MockClient{
		ListCommentsFn: func(ctx context.Context, entityID int, take int, include []string) ([]entity.Comment, error) {
			assert.Equal(t, 200, entityID)
			assert.Equal(t, 50, take)
			assert.Equal(t, []string{"ID", "Description", "Owner"}, include)
			return expectedComments, nil
		},
	}

	tool := NewListCommentsTool(mockClient)
	result := tool.Callback(map[string]interface{}{
		"entityId": float64(200),
		"take":     float64(50),
		"include":  []interface{}{"ID", "Description", "Owner"},
	})

	assert.NotNil(t, result)
	assert.Nil(t, result.IsError)
}
