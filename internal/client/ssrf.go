package client

import (
	"net/url"
	"tp-mcp-go/internal/domain/errors"
)

// validateURL validates that a target URL matches the expected base URL's hostname.
// The hostname check is the core SSRF protection — it ensures requests only go to
// the configured TP domain. Path is not restricted because the TP API uses various
// path prefixes (e.g., /api/v1/, /Upload/, /Attachment/).
func validateURL(targetURL, baseURL string) error {
	if targetURL == "" {
		return &errors.SSRFError{URL: targetURL, Reason: "empty URL"}
	}

	target, err := url.Parse(targetURL)
	if err != nil {
		return &errors.SSRFError{URL: targetURL, Reason: "malformed URL"}
	}

	// Reject relative URLs (must have scheme)
	if target.Scheme == "" || target.Host == "" {
		return &errors.SSRFError{URL: targetURL, Reason: "relative URLs not allowed"}
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return &errors.SSRFError{URL: targetURL, Reason: "malformed base URL"}
	}

	// Check hostname matches exactly (including port)
	if target.Host != base.Host {
		return &errors.SSRFError{URL: targetURL, Reason: "hostname does not match configured domain"}
	}

	return nil
}
