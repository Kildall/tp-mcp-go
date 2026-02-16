package client

import (
	"testing"
	"tp-mcp-go/internal/domain/errors"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		targetURL string
		baseURL   string
		wantErr   bool
		errReason string
	}{
		{
			name:      "valid URL matching domain with /api/v1/ path",
			targetURL: "https://example.com/api/v1/entities",
			baseURL:   "https://example.com/api/v1",
			wantErr:   false,
		},
		{
			name:      "valid URL matching domain with different path",
			targetURL: "https://example.com/Upload/file.pdf",
			baseURL:   "https://example.com/api/v1",
			wantErr:   false,
		},
		{
			name:      "valid URL matching domain with no path",
			targetURL: "https://example.com",
			baseURL:   "https://example.com/api/v1",
			wantErr:   false,
		},
		{
			name:      "mismatched hostname",
			targetURL: "https://evil.com/api/v1/entities",
			baseURL:   "https://example.com/api/v1",
			wantErr:   true,
			errReason: "hostname does not match configured domain",
		},
		{
			name:      "malformed URL",
			targetURL: "ht!tp://invalid url",
			baseURL:   "https://example.com/api/v1",
			wantErr:   true,
			errReason: "malformed URL",
		},
		{
			name:      "empty URL",
			targetURL: "",
			baseURL:   "https://example.com/api/v1",
			wantErr:   true,
			errReason: "empty URL",
		},
		{
			name:      "URL with different port",
			targetURL: "https://example.com:8080/api/v1/entities",
			baseURL:   "https://example.com/api/v1",
			wantErr:   true,
			errReason: "hostname does not match configured domain",
		},
		{
			name:      "relative URL (no scheme)",
			targetURL: "/api/v1/entities",
			baseURL:   "https://example.com/api/v1",
			wantErr:   true,
			errReason: "relative URLs not allowed",
		},
		{
			name:      "URL with matching port in both",
			targetURL: "https://example.com:443/api/v1/entities",
			baseURL:   "https://example.com:443/api/v1",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.targetURL, tt.baseURL)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateURL() expected error but got nil")
					return
				}
				ssrfErr, ok := err.(*errors.SSRFError)
				if !ok {
					t.Errorf("validateURL() expected SSRFError but got %T", err)
					return
				}
				if ssrfErr.Reason != tt.errReason {
					t.Errorf("validateURL() error reason = %q, want %q", ssrfErr.Reason, tt.errReason)
				}
			} else {
				if err != nil {
					t.Errorf("validateURL() unexpected error: %v", err)
				}
			}
		})
	}
}
