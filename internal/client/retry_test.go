package client

import (
	stderrors "errors"
	"fmt"
	"testing"
	"tp-mcp-go/internal/domain/errors"
)

// TestSuccessOnFirstTry verifies that the operation succeeds on the first attempt
func TestSuccessOnFirstTry(t *testing.T) {
	callCount := 0
	operation := func() (string, error) {
		callCount++
		return "success", nil
	}

	result, err := executeWithRetry(operation, 3, 0, 2.0)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if result != "success" {
		t.Errorf("expected 'success', got: %s", result)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got: %d", callCount)
	}
}

// TestSuccessOnSecondAttempt verifies retry on retryable error
func TestSuccessOnSecondAttempt(t *testing.T) {
	callCount := 0
	operation := func() (string, error) {
		callCount++
		if callCount == 1 {
			return "", &errors.APIError{StatusCode: 500, Message: "server error"}
		}
		return "success", nil
	}

	result, err := executeWithRetry(operation, 3, 0, 2.0)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if result != "success" {
		t.Errorf("expected 'success', got: %s", result)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got: %d", callCount)
	}
}

// TestNoRetryOn400 verifies that 400 errors are not retried
func TestNoRetryOn400(t *testing.T) {
	callCount := 0
	operation := func() (string, error) {
		callCount++
		return "", &errors.APIError{StatusCode: 400, Message: "bad request"}
	}

	result, err := executeWithRetry(operation, 3, 0, 2.0)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if result != "" {
		t.Errorf("expected empty result, got: %s", result)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (no retry), got: %d", callCount)
	}

	var apiErr *errors.APIError
	if !stderrors.As(err, &apiErr) || apiErr.StatusCode != 400 {
		t.Errorf("expected APIError with status 400, got: %v", err)
	}
}

// TestNoRetryOn401 verifies that 401 errors are not retried
func TestNoRetryOn401(t *testing.T) {
	callCount := 0
	operation := func() (string, error) {
		callCount++
		return "", &errors.APIError{StatusCode: 401, Message: "unauthorized"}
	}

	result, err := executeWithRetry(operation, 3, 0, 2.0)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if result != "" {
		t.Errorf("expected empty result, got: %s", result)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (no retry), got: %d", callCount)
	}

	var apiErr *errors.APIError
	if !stderrors.As(err, &apiErr) || apiErr.StatusCode != 401 {
		t.Errorf("expected APIError with status 401, got: %v", err)
	}
}

// TestExhaustedRetriesReturnsLastError verifies that after maxRetries, the last error is returned
func TestExhaustedRetriesReturnsLastError(t *testing.T) {
	callCount := 0
	operation := func() (string, error) {
		callCount++
		return "", &errors.APIError{StatusCode: 503, Message: fmt.Sprintf("attempt %d", callCount)}
	}

	result, err := executeWithRetry(operation, 2, 0, 2.0)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if result != "" {
		t.Errorf("expected empty result, got: %s", result)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls (1 initial + 2 retries), got: %d", callCount)
	}

	var apiErr *errors.APIError
	if !stderrors.As(err, &apiErr) || apiErr.Message != "attempt 3" {
		t.Errorf("expected last error with 'attempt 3', got: %v", err)
	}
}

// TestAttemptCount verifies the correct number of attempts
func TestAttemptCount(t *testing.T) {
	tests := []struct {
		name       string
		maxRetries int
		wantCalls  int
	}{
		{"zero retries", 0, 1},
		{"one retry", 1, 2},
		{"three retries", 3, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			operation := func() (int, error) {
				callCount++
				return 0, &errors.APIError{StatusCode: 500, Message: "error"}
			}

			executeWithRetry(operation, tt.maxRetries, 0, 2.0)

			if callCount != tt.wantCalls {
				t.Errorf("expected %d calls, got: %d", tt.wantCalls, callCount)
			}
		})
	}
}
