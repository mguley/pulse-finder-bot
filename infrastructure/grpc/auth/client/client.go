package client

import (
	"context"
	"errors"
	"fmt"
	authv1 "infrastructure/proto/auth/gen"

	"google.golang.org/grpc"
)

// AuthClient is a high-level wrapper over the underlying gRPC client connection to AuthService.
// It manages the gRPC connection and provides methods to perform GenerateToken calls.
type AuthClient struct {
	conn   *grpc.ClientConn         // The underlying gRPC client connection.
	client authv1.AuthServiceClient // The generated AuthService client.
}

// NewAuthClient creates a new instance of AuthClient based on the provided configuration.
func NewAuthClient(env, address string) (*AuthClient, error) {
	var (
		conn   *grpc.ClientConn
		config *Config
		err    error
	)

	switch env {
	case "prod":
		// Production environment with TLS (certificate file is optional here, adjust as needed).
		conn, config, err = NewGRPCClient(
			WithServerAddress(address),
			WithTLS(""))
	case "dev":
		// Development environment without TLS.
		conn, config, err = NewGRPCClient(
			WithServerAddress(address))
	default:
		return nil, errors.New("unsupported environment; must be \"prod\" or \"dev\"")
	}

	if err != nil {
		return nil, fmt.Errorf("create gRPC client: %w", err)
	}
	fmt.Printf("Connected to server: %s (TLS enabled: %v)\n", address, config.TLSEnabled)
	return &AuthClient{
		conn:   conn,
		client: authv1.NewAuthServiceClient(conn),
	}, nil
}

// GenerateToken calls the AuthService's GenerateToken RPC method.
// It returns the generated token or an error if the call failed.
func (c *AuthClient) GenerateToken(ctx context.Context, issuer string, scopes []string) (string, error) {
	resp, err := c.client.GenerateToken(ctx, &authv1.GenerateTokenRequest{
		Issuer: issuer,
		Scopes: scopes,
	})
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return resp.GetToken(), nil
}

// Close closes the underlying gRPC connection.
func (c *AuthClient) Close() error {
	return c.conn.Close()
}
