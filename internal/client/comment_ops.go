package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"tp-mcp-go/internal/domain/entity"
)

// CreateComment creates a private comment on an entity
func (c *httpClient) CreateComment(ctx context.Context, entityID int, description string) (*entity.Comment, error) {
	url := fmt.Sprintf("%s/Comments", c.baseURL)
	body := map[string]any{
		"Description": description,
		"General":     map[string]any{"Id": entityID},
		"IsPrivate":   true,
	}
	data, err := c.doPost(ctx, url, body)
	if err != nil {
		return nil, err
	}
	var comment entity.Comment
	if err := json.Unmarshal(data, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListComments lists comments for an entity
func (c *httpClient) ListComments(ctx context.Context, entityID int, take int, include []string) ([]entity.Comment, error) {
	if len(include) == 0 {
		include = []string{"Description", "CreateDate", "Owner"}
	}

	url := fmt.Sprintf("%s/Comments?where=General.Id eq %d&take=%d&include=[%s]&orderBy=CreateDate desc",
		c.baseURL, entityID, take, strings.Join(include, ","))

	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}

	var apiResp entity.APIResponse
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, err
	}

	// Convert Items to Comment structs
	var comments []entity.Comment
	for _, item := range apiResp.Items {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			continue
		}
		var comment entity.Comment
		if err := json.Unmarshal(itemBytes, &comment); err != nil {
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
