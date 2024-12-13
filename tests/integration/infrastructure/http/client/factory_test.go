package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFactory_CreateDefaultClient tests the creation of a default HTTP client.
func TestFactory_CreateDefaultClient(t *testing.T) {
	container := SetupTestContainer()
	factory := container.HttpFactory.Get()

	// Create a default client
	client := factory.CreateDefaultClient()
	require.NotNil(t, client, "Default HTTP client creation failed")

	// Validate the User-Agent is set
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org?format=json", http.NoBody)
	require.NoError(t, err, "Failed to create HTTP request")

	resp, err := client.Do(req)
	assert.NotNil(t, resp, "Response should not be nil")
	assert.NoError(t, err, "Failed to execute HTTP request")
	assert.Contains(t, req.Header.Get("User-Agent"), "AppleWebKit/537.36 (KHTML, like Gecko) Chrome")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestFactory_CreateSocks5Client tests the creation of a SOCKS5 HTTP client.
func TestFactory_CreateSocks5Client(t *testing.T) {
	container := SetupTestContainer()
	factory := container.HttpFactory.Get()
	config := container.Config.Get()

	// Settings
	h := config.Proxy.Host
	p := config.Proxy.Port
	timeout := 10 * time.Second

	// Create a SOCKS5 client
	client, err := factory.CreateSocks5Client(h, p, timeout)
	require.NoError(t, err, "Failed to create Socks5 client")

	// Validate the User-Agent is set
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org?format=json", http.NoBody)
	require.NoError(t, err, "Failed to create HTTP request")

	resp, err := client.Do(req)
	assert.NotNil(t, resp, "Response should not be nil")
	assert.NoError(t, err, "Failed to execute HTTP request")
	assert.Contains(t, req.Header.Get("User-Agent"), "AppleWebKit/537.36 (KHTML, like Gecko) Chrome")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
