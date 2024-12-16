package sitemap

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// SetupTestContainer initializes the TestContainer and handles cleanup.
func SetupTestContainer(t *testing.T) *TestContainer {
	c := NewTestContainer()
	config := c.Config.Get()

	// Cleanup resources after tests
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Drop the test database to clean up after tests
		err := c.MongoClient.Get().Database(config.Mongo.DB).Drop(ctx)
		if err != nil {
			fmt.Printf("failed to drop test database: %v\n", err)
		}

		// Disconnect the MongoDB client
		err = c.MongoClient.Get().Disconnect(ctx)
		if err != nil {
			fmt.Printf("failed to disconnect test database: %v\n", err)
		}
	})

	return c
}
