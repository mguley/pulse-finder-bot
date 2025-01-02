package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds the client configuration settings.
type Config struct {
	TLSEnabled    bool   // Whether to use TLS
	ServerAddress string // Target server address
	CertFile      string // Path to the certificate file (for TLS)
}

// Option defines a functional option for configuring the client.
type Option func(*Config)

// WithTLS enables TLS and sets the certificate file for the client.
func WithTLS(certFile string) Option {
	return func(c *Config) {
		c.TLSEnabled = true
		c.CertFile = certFile
	}
}

// WithServerAddress sets the target server address.
func WithServerAddress(address string) Option {
	return func(c *Config) {
		c.ServerAddress = address
	}
}

// NewGRPCClient initializes a gRPC client connection with the provided options.
func NewGRPCClient(opts ...Option) (*grpc.ClientConn, *Config, error) {
	config := &Config{
		TLSEnabled: false,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	// Set up connection options
	var dialOpts []grpc.DialOption
	if config.TLSEnabled {
		cred, err := getTransportCredentials(config.CertFile)
		if err != nil {
			return nil, nil, fmt.Errorf("could not load TLS credentials: %w", err)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(cred))
	} else {
		// Use insecure credentials for local development
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Create the gRPC client connection
	conn, err := grpc.NewClient(config.ServerAddress, dialOpts...)
	if err != nil {
		return nil, nil, err
	}
	return conn, config, nil
}

// getTransportCredentials determines the correct transport credentials based on the certFile.
func getTransportCredentials(certFile string) (credentials.TransportCredentials, error) {
	if certFile != "" {
		// Use client-side TLS with the provided certificate
		return credentials.NewClientTLSFromFile(certFile, "")
	}
	// Use system CA trust store for validation
	return credentials.NewTLS(nil), nil
}
