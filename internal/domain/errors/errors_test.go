package errors

import (
	"fmt"
	"testing"
)

func TestInvalidEntityTypeError(t *testing.T) {
	err := &InvalidEntityTypeError{Type: "UnknownType"}
	expected := "invalid entity type: UnknownType"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{Field: "name", Message: "cannot be empty"}
	expected := "validation error on field name: cannot be empty"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestAPIError(t *testing.T) {
	t.Run("without context", func(t *testing.T) {
		err := &APIError{StatusCode: 404, Message: "not found"}
		expected := "API error 404: not found"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("with context", func(t *testing.T) {
		err := &APIError{StatusCode: 500, Message: "server error", Context: "fetching entity"}
		expected := "API error 500: server error (context: fetching entity)"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})
}

func TestSSRFError(t *testing.T) {
	err := &SSRFError{URL: "http://localhost", Reason: "private IP"}
	expected := "SSRF validation failed for URL http://localhost: private IP"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "400 not retryable",
			err:      &APIError{StatusCode: 400},
			expected: false,
		},
		{
			name:     "401 not retryable",
			err:      &APIError{StatusCode: 401},
			expected: false,
		},
		{
			name:     "500 retryable",
			err:      &APIError{StatusCode: 500},
			expected: true,
		},
		{
			name:     "non-API error retryable",
			err:      &ValidationError{Field: "test", Message: "error"},
			expected: true,
		},
		{
			name:     "wrapped 400 not retryable",
			err:      fmt.Errorf("wrapped: %w", &APIError{StatusCode: 400}),
			expected: false,
		},
		{
			name:     "wrapped 500 retryable",
			err:      fmt.Errorf("wrapped: %w", &APIError{StatusCode: 500}),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryable(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		token    string
		expected string
	}{
		{
			name:     "masks token in message",
			msg:      "error: token abc123 is invalid",
			token:    "abc123",
			expected: "error: token *** is invalid",
		},
		{
			name:     "handles empty token",
			msg:      "error: token abc123 is invalid",
			token:    "",
			expected: "error: token abc123 is invalid",
		},
		{
			name:     "handles message without token",
			msg:      "error: no token here",
			token:    "abc123",
			expected: "error: no token here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskToken(tt.msg, tt.token)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
