package strategies

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// RetryStrategy defines the contract for a retry strategy.
// A retry strategy determines how long to wait before retrying an operation after failure.
type RetryStrategy interface {
	// WaitDuration calculates the duration to wait before the next retry attempt.
	// Returns an error if the attempt exceeds a configured maximum.
	WaitDuration(attempt int) (time.Duration, error)
}

// ExponentialBackoffStrategy implements an exponential backoff retry strategy.
// It increases the wait time exponentially with each retry attempt, up to a maximum delay.
type ExponentialBackoffStrategy struct {
	baseDelay   time.Duration // Minimum delay between retries (e.g., 100ms).
	maxDelay    time.Duration // Maximum delay between retries (e.g., 5s).
	maxAttempts int           // Maximum number of retry attempts (0 for unlimited).
	multiplier  float64       // Growth multiplier for exponential backoff (e.g., 2.0)
}

// NewExponentialBackoffStrategy initializes a new ExponentialBackoffStrategy with provided parameters.
func NewExponentialBackoffStrategy(baseDelay, maxDelay time.Duration, attempts int,
	multiplier float64) *ExponentialBackoffStrategy {
	// Apply default values where parameters are invalid or zero
	if baseDelay <= 0 {
		baseDelay = 5 * time.Second
	}
	if maxDelay <= 0 {
		maxDelay = 15 * time.Second
	}
	if multiplier <= 0 {
		multiplier = 2.0
	}
	return &ExponentialBackoffStrategy{
		baseDelay:   baseDelay,
		maxDelay:    maxDelay,
		maxAttempts: attempts,
		multiplier:  multiplier,
	}
}

// WaitDuration calculates the duration to wait before the next retry attempt.
// It uses an exponential backoff formula: BaseDelay * Multiplier^attempt.
func (s *ExponentialBackoffStrategy) WaitDuration(attempt int) (time.Duration, error) {
	if err := s.validate(attempt); err != nil {
		return 0, fmt.Errorf("exponential backoff strategy: %w", err)
	}

	// Calculate exponential delay
	delay := float64(s.baseDelay) * math.Pow(s.multiplier, float64(attempt))
	if delay > float64(s.maxDelay) {
		delay = float64(s.maxDelay)
	}

	return time.Duration(delay), nil
}

// validate checks the validity of the ExponentialBackoffStrategy configuration and the given attempt number.
func (s *ExponentialBackoffStrategy) validate(a int) error {
	if a < 0 {
		return errors.New("attempts must be greater than zero")
	}
	if s.maxAttempts > 0 && a >= s.maxAttempts {
		return errors.New("maximum retry attempts exceeded")
	}
	return nil
}
