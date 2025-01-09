package cron

import (
	"context"
	"fmt"
	"testing"
	authClient "tests/integration/application/scheduler/cron/auth/client"
	vacancyClient "tests/integration/application/scheduler/cron/vacancy/client"
	"time"
)

type TestEnvironment struct {
	AuthClientSetup    *authClient.TestEnvironment
	VacancyClientSetup *vacancyClient.TestEnvironment
	TestContainer      *TestContainer
}

// SetupTestEnvironment initializes the TestContainer and handles cleanup.
func SetupTestEnvironment(t *testing.T) *TestContainer {
	c := &TestEnvironment{
		AuthClientSetup:    authClient.SetupTestEnvironment(t),
		VacancyClientSetup: vacancyClient.SetupTestEnvironment(t),
		TestContainer:      NewTestContainer(),
	}

	// Cleanup resources after tests
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Drop the test database to clean up after tests
		err := c.TestContainer.MongoClient.Get().Database(c.TestContainer.Config.Get().Mongo.DB).Drop(ctx)
		if err != nil {
			fmt.Printf("failed to drop test database: %v\n", err)
		}

		// Disconnect the MongoDB client
		err = c.TestContainer.MongoClient.Get().Disconnect(ctx)
		if err != nil {
			fmt.Printf("failed to disconnect test database: %v\n", err)
		}
	})

	return c.TestContainer
}
