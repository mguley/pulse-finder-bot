package client

import (
	"application/dependency"
	"infrastructure/grpc/vacancy/client"
	"log"
	"tests/integration/infrastructure/grpc/vacancy/client/server"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	MockVacancyServiceServer dependency.LazyDependency[*server.MockVacancyService]
	TestServerContainer      dependency.LazyDependency[*server.TestServerContainer]
	VacancyClient            dependency.LazyDependency[*client.VacancyClient]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.MockVacancyServiceServer = dependency.LazyDependency[*server.MockVacancyService]{
		InitFunc: server.NewMockVacancyService,
	}
	c.TestServerContainer = dependency.LazyDependency[*server.TestServerContainer]{
		InitFunc: func() *server.TestServerContainer {
			grpcServer, err := server.NewTestServerContainer(c.MockVacancyServiceServer.Get())
			if err != nil {
				log.Fatalf("Failed to create gRPC test server: %v", err)
			}
			return grpcServer
		},
	}
	c.VacancyClient = dependency.LazyDependency[*client.VacancyClient]{
		InitFunc: func() *client.VacancyClient {
			grpcClient, err := client.NewVacancyClient("dev", c.TestServerContainer.Get().Address)
			if err != nil {
				log.Fatalf("Failed to create gRPC test client: %v", err)
			}
			return grpcClient
		},
	}

	return c
}
