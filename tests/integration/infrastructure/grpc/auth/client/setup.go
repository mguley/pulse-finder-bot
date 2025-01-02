package client

import (
	"infrastructure/grpc/auth/client"
	"testing"
	"tests/integration/infrastructure/grpc/auth/client/server"
)

// TestEnvironment encapsulates the mock server and the gRPC client for integration tests.
type TestEnvironment struct {
	Server *server.TestServerContainer // The mock gRPC server container.
	Client *client.AuthClient          // The gRPC client connected to the mock server.
}

// SetupTestEnvironment initializes the test environment with a mock server and client.
//
// It starts a mock gRPC server and configures a client to communicate with it.
// Resources are automatically cleaned up after tests using t.Cleanup.
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	container := NewTestContainer()

	grpcServer := container.TestServerContainer.Get()
	grpcClient := container.AuthClient.Get()

	// Ensure resources are cleaned up after tests.
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
