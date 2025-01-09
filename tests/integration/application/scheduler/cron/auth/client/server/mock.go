package server

import (
	"context"
	authv1 "infrastructure/proto/auth/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockAuthService is a mock implementation of AuthServiceServer for testing.
type MockAuthService struct {
	authv1.UnimplementedAuthServiceServer // Ensures forward compatibility with the gRPC interface.
}

// NewMockAuthService creates and returns a new instance of MockAuthService.
func NewMockAuthService() *MockAuthService { return &MockAuthService{} }

// GenerateToken simulates the behavior of the GenerateToken RPC method.
// It validates the request and generates a mock JWT token or returns an error.
func (s *MockAuthService) GenerateToken(ctx context.Context, req *authv1.GenerateTokenRequest) (*authv1.GenerateTokenResponse, error) {
	if req.GetIssuer() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "issuer (iss) must not be empty")
	}
	if len(req.GetScopes()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "scope must not be empty")
	}

	// Simulate generating a token
	token := "mock-jwt-token"
	return &authv1.GenerateTokenResponse{Token: token}, nil
}
