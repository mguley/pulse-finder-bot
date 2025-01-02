package client

import (
	"application/dependency"
	"infrastructure/grpc/auth/client"
	"log"
	"tests/integration/infrastructure/grpc/auth/client/server"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	MockAuthServiceServer dependency.LazyDependency[*server.MockAuthService]
	TestServerContainer   dependency.LazyDependency[*server.TestServerContainer]
	AuthClient            dependency.LazyDependency[*client.AuthClient]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.MockAuthServiceServer = dependency.LazyDependency[*server.MockAuthService]{
		InitFunc: server.NewMockAuthService,
	}
	c.TestServerContainer = dependency.LazyDependency[*server.TestServerContainer]{
		InitFunc: func() *server.TestServerContainer {
			grpcServer, err := server.NewTestServerContainer(c.MockAuthServiceServer.Get())
			if err != nil {
				log.Fatalf("Failed to create gRPC test server: %v", err)
			}
			return grpcServer
		},
	}
	c.AuthClient = dependency.LazyDependency[*client.AuthClient]{
		InitFunc: func() *client.AuthClient {
			grpcClient, err := client.NewAuthClient("dev", c.TestServerContainer.Get().Address)
			if err != nil {
				log.Fatalf("Failed to create gRPC test client: %v", err)
			}
			return grpcClient
		},
	}

	return c
}
