package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"tp-mcp-go/internal/domain/entity"
)

// ListAttachments lists attachments for an entity
func (c *httpClient) ListAttachments(ctx context.Context, entityID int, take int) ([]entity.Attachment, error) {
	url := fmt.Sprintf("%s/Attachments?where=General.Id eq %d&take=%d",
		c.baseURL, entityID, take)

	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}

	var apiResp entity.APIResponse
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, err
	}

	var attachments []entity.Attachment
	for _, item := range apiResp.Items {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			continue
		}
		var att entity.Attachment
		if err := json.Unmarshal(itemBytes, &att); err != nil {
			continue
		}
		attachments = append(attachments, att)
	}
	return attachments, nil
}

// GetAttachmentMetadata retrieves metadata for a single attachment
func (c *httpClient) GetAttachmentMetadata(ctx context.Context, attachmentID int) (*entity.Attachment, error) {
	url := fmt.Sprintf("%s/Attachments/%d", c.baseURL, attachmentID)
	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	var att entity.Attachment
	if err := json.Unmarshal(data, &att); err != nil {
		return nil, err
	}
	return &att, nil
}

// DownloadAttachment downloads attachment content, validating the URL first
func (c *httpClient) DownloadAttachment(ctx context.Context, uri string) ([]byte, string, error) {
	// Build full URL if relative
	downloadURL := uri
	if !strings.HasPrefix(uri, "http") {
		downloadURL = c.baseURL[:strings.Index(c.baseURL, "/api/v1")] + uri
	}

	// SSRF validation
	if err := validateURL(downloadURL, c.baseURL); err != nil {
		return nil, "", err
	}

	// Use a dedicated HTTP client that re-applies auth on redirects.
	// Go's default client drops headers/query params on redirect, which
	// causes the TP server to return a login page instead of the file.
	downloadClient := &http.Client{
		Timeout: c.httpClient.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			// Re-apply auth on every redirect so the access token isn't lost
			c.auth.ApplyAuth(req)
			return nil
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, "", err
	}
	c.auth.ApplyAuth(req)

	resp, err := downloadClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(body))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	mimeType := resp.Header.Get("Content-Type")
	if strings.Contains(mimeType, "text/html") {
		preview := string(content)
		if len(preview) > 500 {
			preview = preview[:500]
		}
		return nil, "", fmt.Errorf("download returned HTML instead of file content (likely authentication failure): %s", preview)
	}
	return content, mimeType, nil
}
