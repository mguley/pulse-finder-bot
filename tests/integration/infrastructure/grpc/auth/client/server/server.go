package server

import (
	"fmt"
	authv1 "infrastructure/proto/auth/gen"
	"log"
	"net"

	"google.golang.org/grpc"
)

// TestServerContainer manages the lifecycle of a test gRPC server.
type TestServerContainer struct {
	grpcServer *grpc.Server // The gRPC server instance.
	listener   net.Listener // The listener for incoming gRPC connections.
	Address    string       // The address the server is listening on.
}

// NewTestServerContainer initializes and starts a test gRPC server.
// It registers the provided AuthServiceServer implementation.
func NewTestServerContainer(authServer authv1.AuthServiceServer) (*TestServerContainer, error) {
	// Create a listener on a random available port.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authServer)

	// Start the server in a separate goroutine.
	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return &TestServerContainer{
		grpcServer: grpcServer,
		listener:   listener,
		Address:    listener.Addr().String(),
	}, nil
}

// Stop gracefully stops the test gRPC server.
func (s *TestServerContainer) Stop() {
	s.grpcServer.GracefulStop()
}
