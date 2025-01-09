package client

import (
	"infrastructure/grpc/vacancy/client"
	"testing"
	"tests/integration/application/scheduler/cron/vacancy/client/server"
)

// TestEnvironment encapsulates the mock server and the gRPC client for integration tests.
type TestEnvironment struct {
	Server *server.TestServerContainer // The mock gRPC server container.
	Client *client.VacancyClient       // The gRPC client connected to the mock server.
}

// SetupTestEnvironment initializes the test environment with a mock server and client.
//
// It starts a mock gRPC server and configures a client to communicate with it.
// Resources are automatically cleaned up after tests using t.Cleanup.
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	container := NewTestContainer()

	grpcServer := container.TestServerContainer.Get()
	grpcClient := container.VacancyClient.Get()

	t.Cleanup(func() {
		grpcServer.Stop()
		if err := grpcClient.Close(); err != nil {
			t.Logf("failed to close gRPC client: %v", err)
		}
	})

	return &TestEnvironment{
		Server: grpcServer,
		Client: grpcClient,
	}
}
