package commands

import (
	"log"
	"testing"
)

// SetupTestContainer initializes the TestContainer and handles cleanup.
func SetupTestContainer(t *testing.T) *TestContainer {
	c := NewTestContainer()

	t.Cleanup(func() {
		err := c.ProxyConnection.Get().Close()
		if err != nil {
			log.Printf("failed to close proxy connection: %v", err)
		}
	})

	return c
}
