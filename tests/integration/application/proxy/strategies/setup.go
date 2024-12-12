package strategies

import "testing"

// SetupTestContainer initializes the TestContainer.
func SetupTestContainer(t *testing.T) *TestContainer {
	return NewTestContainer()
}
