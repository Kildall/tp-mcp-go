package errors

import (
	stderrors "errors"
	"fmt"
	"strings"
)

// InvalidEntityTypeError is returned when an invalid entity type is provided
type InvalidEntityTypeError struct {
	Type string
}

func (e *InvalidEntityTypeError) Error() string {
	return fmt.Sprintf("invalid entity type: %s", e.Type)
}

// ValidationError is returned for validation failures
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// APIError is returned when the TP API returns a non-success response
type APIError struct {
	StatusCode int
	Message    string // Parsed/clean error message
	RawBody    string // Full response body (with token masked)
	Context    string // Method + URL
}

func (e *APIError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("API error %d: %s (context: %s)", e.StatusCode, e.Message, e.Context)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// ParseTPErrorBody extracts a human-readable message from a TP API error response.
// The TP API returns XML error bodies like:
//
//	<Error><Status>BadRequest</Status><Message>Error during parameters parsing.</Message>...</Error>
func ParseTPErrorBody(body string) string {
	start := strings.Index(body, "<Message>")
	end := strings.Index(body, "</Message>")
	if start != -1 && end != -1 && end > start {
		return body[start+len("<Message>") : end]
	}
	// If no XML message found, return body as-is (trimmed)
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return "empty response"
	}
	return trimmed
}

// SSRFError is returned when URL validation fails
type SSRFError struct {
	URL    string
	Reason string
}

func (e *SSRFError) Error() string {
	return fmt.Sprintf("SSRF validation failed for URL %s: %s", e.URL, e.Reason)
}

// IsRetryable returns whether the error should be retried
func IsRetryable(err error) bool {
	var apiErr *APIError
	if stderrors.As(err, &apiErr) {
		return apiErr.StatusCode != 400 && apiErr.StatusCode != 401
	}
	return true // non-API errors are retryable by default
}

// MaskToken replaces access token occurrences in error messages
func MaskToken(msg, token string) string {
	if token == "" {
		return msg
	}
	return strings.ReplaceAll(msg, token, "***")
}
