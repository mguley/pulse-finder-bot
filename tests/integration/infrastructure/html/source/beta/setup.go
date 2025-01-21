package beta

import "testing"

// SetupTestContainer initializes the TestContainer and handles cleanup.
func SetupTestContainer(t *testing.T) *TestContainer {
	return NewTestContainer()
}
