package commands

import (
	"context"
	"errors"
	"fmt"
	httpClient "infrastructure/http/client"
	"infrastructure/proxy/client"
	"net/http"
	"time"
)

// StatusCommand provides functionality to check the status of a SOCKS5 proxy.
type StatusCommand struct {
	factory      *httpClient.Factory  // Factory for creating HTTP clients.
	socks5Client *client.Socks5Client // SOCKS5 client for managing proxy interactions.
	host         string               // Hostname or IP address of the SOCKS5 proxy.
	port         string               // Port number of the SOCKS5 proxy.
	timeout      time.Duration        // Timeout duration for proxy interactions.
}

// NewStatusCommand creates a new instance of StatusCommand.
func NewStatusCommand(h, p string, c *client.Socks5Client, f *httpClient.Factory, t time.Duration) *StatusCommand {
	return &StatusCommand{host: h, port: p, factory: f, socks5Client: c, timeout: t}
}

// Execute performs a status check by connecting to the specified URL through the SOCKS5 proxy.
func (cmd *StatusCommand) Execute(url string) (string, error) {
	// Validate inputs.
	if cmd.host == "" || cmd.port == "" {
		return "", errors.New("proxy host or port is not configured")
	}

	// Create the HTTP client.
	c, err := cmd.createSocks5Client()
	if err != nil {
		return "", fmt.Errorf("create socks5 client: %w", err)
	}

	// Ping the specified URL.
	return cmd.pingProxy(c, url)
}

// createSocks5Client initializes an HTTP client configured for SOCKS5 proxy usage.
func (cmd *StatusCommand) createSocks5Client() (*http.Client, error) {
	return cmd.factory.CreateSocks5Client(cmd.host, cmd.port, cmd.timeout)
}

// pingProxy performs HTTP request to URL through the SOCKS5 proxy to check connectivity.
func (cmd *StatusCommand) pingProxy(c *http.Client, url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.timeout)
	defer cancel()

	// Ping the URL using the SOCKS5 client.
	res, err := cmd.socks5Client.Ping(ctx, c, url)
	if err != nil {
		return "", fmt.Errorf("ping proxy for URL %s: %w", url, err)
	}
	return res, nil
}
