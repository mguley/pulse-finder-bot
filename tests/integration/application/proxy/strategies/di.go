package strategies

import (
	"application/dependency"
	"application/proxy/strategies"
	"time"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	RetryStrategy dependency.LazyDependency[strategies.RetryStrategy]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.RetryStrategy = dependency.LazyDependency[strategies.RetryStrategy]{
		InitFunc: func() strategies.RetryStrategy {
			baseDelay := 5 * time.Second
			maxDelay := 30 * time.Second
			maxAttempts := 5
			multiplier := 2.0
			return strategies.NewExponentialBackoffStrategy(baseDelay, maxDelay, maxAttempts, multiplier)
		},
	}

	return c
}
