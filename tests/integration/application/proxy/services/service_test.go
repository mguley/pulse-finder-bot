package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestService_HttpClient_Success tests the creation and reuse of an HTTP client.
func TestService_HttpClient_Success(t *testing.T) {
	container := SetupTestContainer(t)
	service := container.ProxyService.Get()

	// Retrieve the HTTP client
	c, err := service.HttpClient()
	require.NoError(t, err, "Failed to create HTTP client")
	require.NotNil(t, c, "Expected HTTP client to be non-nil")

	// Validate that subsequent calls return the same client
	c2, err := service.HttpClient()
	require.NoError(t, err, "Failed to retrieve HTTP client")
	assert.Equal(t, c, c2, "Expected the same HTTP client instance")
}

// TestService_HttpClient_Failure tests the behavior when an invalid proxy is used.
func TestService_HttpClient_Failure(t *testing.T) {
	container := SetupTestContainer(t)

	// Override the proxy host and port with invalid values
	container.Config.Get().Proxy.Host = "invalid_host"
	container.Config.Get().Proxy.Port = "9999"

	service := container.ProxyService.Get()

	// Retrieve the HTTP client
	c, err := service.HttpClient()
	require.NoError(t, err, "HTTP client creation should not fail before making requests")
	require.NotNil(t, c, "Expected HTTP client to be non-nil")

	// Attempt to make an HTTP request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.ipify.org?format=json", http.NoBody)
	require.NoError(t, err, "Failed to create HTTP request")

	// Execute the request
	res, err := c.Do(req)
	require.Error(t, err, "Expected an error due to invalid proxy configuration")
	assert.Contains(t, err.Error(), "dial", "Error message should indicate a connection failure")
	assert.Nil(t, res, "Expected no response from HTTP request")
}
