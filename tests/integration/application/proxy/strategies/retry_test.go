package strategies

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRetryStrategy_Success verifies that the exponential backoff strategy calculates the correct delays
// for various retry attempts. It ensures delays grow exponentially and are capped at the maximum delay.
func TestRetryStrategy_Success(t *testing.T) {
	container := SetupTestContainer()
	strategy := container.RetryStrategy.Get()

	// Test the wait duration calculation for various attempts
	var delays []time.Duration
	for i := 0; i < 5; i++ {
		delay, err := strategy.WaitDuration(i)
		require.NoError(t, err, "Failed to wait for strategy #%d", i)
		delays = append(delays, delay)
	}

	// Assert that delays grow exponentially and are capped at maxDelay
	for i := 1; i < len(delays); i++ {
		assert.GreaterOrEqual(t, delays[i], delays[i-1], "Delays should grow exponentially")
	}
	assert.LessOrEqual(t, delays[len(delays)-1], 30*time.Second, "Delay should not exceed maxDelay")
}

// TestRetryStrategy_MaxAttemptsExceeded verifies that the exponential backoff strategy
// returns an error and zero delay when the number of attempts exceeds the configured maximum.
func TestRetryStrategy_MaxAttemptsExceeded(t *testing.T) {
	container := SetupTestContainer()
	strategy := container.RetryStrategy.Get()

	// Test exceeding max attempts
	delay, err := strategy.WaitDuration(5)
	assert.Equal(t, time.Duration(0), delay, "Expected delay to be 0 as max attempts exceeded")
	require.Error(t, err, "Expected an error when max attempts are exceeded")
	assert.Contains(t, err.Error(), "maximum retry attempts exceeded", "Error should indicate max attempts exceeded")
}

// TestRetryStrategy_InvalidAttempt verifies that the exponential backoff strategy
// returns an error and zero delay for invalid (negative) attempt numbers.
func TestRetryStrategy_InvalidAttempt(t *testing.T) {
	container := SetupTestContainer()
	strategy := container.RetryStrategy.Get()

	// Test invalid (negative) attempt
	delay, err := strategy.WaitDuration(-1)
	assert.Equal(t, time.Duration(0), delay, "Expected delay to be 0 for invalid attempt number")
	require.Error(t, err, "Expected an error when attempt number is negative")
	assert.Contains(t, err.Error(), "attempts must be greater than zero", "Error should indicate invalid attempt")
}
