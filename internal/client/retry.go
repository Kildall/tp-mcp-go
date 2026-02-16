package client

import (
	"math"
	"time"
	"tp-mcp-go/internal/domain/errors"
)

// executeWithRetry retries an operation with exponential backoff.
// It does not retry on 400/401 errors (checked via errors.IsRetryable).
func executeWithRetry[T any](operation func() (T, error), maxRetries int, initialDelay time.Duration, backoffFactor float64) (T, error) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, err := operation()
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !errors.IsRetryable(err) {
			var zero T
			return zero, err
		}
		if attempt < maxRetries {
			delay := time.Duration(float64(initialDelay) * math.Pow(backoffFactor, float64(attempt)))
			time.Sleep(delay)
		}
	}
	var zero T
	return zero, lastErr
}
