package commands

import (
	"context"
	"errors"
	"fmt"
	httpClient "infrastructure/http/client"
	"io"
	"net/http"
	"time"
)

// StatusCommand checks the status of a SOCKS5 proxy.
type StatusCommand struct {
	factory *httpClient.Factory // Factory for creating HTTP clients.
	host    string              // Hostname or IP address of the SOCKS5 proxy.
	port    string              // Port number of the SOCKS5 proxy.
	timeout time.Duration       // Timeout duration for proxy interactions.
}

// NewStatusCommand creates a new instance of StatusCommand.
func NewStatusCommand(host, port string, factory *httpClient.Factory, timeout time.Duration) *StatusCommand {
	return &StatusCommand{host: host, port: port, factory: factory, timeout: timeout}
}

// Execute checks the connectivity through the SOCKS5 proxy.
func (c *StatusCommand) Execute(url string) (result string, err error) {
	if c.host == "" || c.port == "" {
		return "", errors.New("proxy host or port is not configured")
	}

	client, err := c.createSocks5Client()
	if err != nil {
		return "", fmt.Errorf("create socks5 client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.Ping(ctx, client, url)
}

// createSocks5Client returns HTTP client configured for SOCKS5 proxy usage.
func (c *StatusCommand) createSocks5Client() (client *http.Client, err error) {
	return c.factory.CreateSocks5Client(c.host, c.port, c.timeout)
}

// Ping checks connectivity through the SOCKS5 proxy by sending HTTP request.
func (c *StatusCommand) Ping(ctx context.Context, client *http.Client, url string) (result string, err error) {
	var (
		request  *http.Request
		response *http.Response
		body     []byte
	)

	if request, err = http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody); err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	if response, err = client.Do(request); err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			fmt.Printf("close response body err:%v", err)
		}
	}()

	if body, err = io.ReadAll(response.Body); err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}
	return string(body), nil
}
