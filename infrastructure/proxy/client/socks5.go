package client

import (
	"context"
	"domain/useragent"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

// Socks5Client provides functionality to interact with a SOCKS5 proxy.
type Socks5Client struct {
	agent useragent.Generator
}

// NewSocks5Client creates a new instance of Socks5Client.
func NewSocks5Client(u useragent.Generator) *Socks5Client {
	return &Socks5Client{agent: u}
}

// HttpClient creates an HTTP client configured to use a SOCKS5 proxy.
func (s *Socks5Client) HttpClient(host, port string, timeout time.Duration) (*http.Client, error) {
	if err := s.validate(host, port); err != nil {
		return nil, err
	}

	proxyURI := fmt.Sprintf("%s:%s", host, port)
	dialer, err := proxy.SOCKS5("tcp", proxyURI, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SOCKS5 proxy: %s", err)
	}

	// Define a context-aware dialer
	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}
	httpTransport := &http.Transport{
		DialContext: dialContext,
	}

	uTransport := &RoundTripWithUserAgent{
		rt:        httpTransport,
		userAgent: s.agent.Generate(),
	}

	httpClient := &http.Client{
		Transport: uTransport,
		Timeout:   timeout,
	}

	return httpClient, nil
}

// Ping verifies connectivity through the SOCKS5 proxy by sending an HTTP request.
func (s *Socks5Client) Ping(ctx context.Context, c *http.Client, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create ping request: %w", err)
	}

	// Execute the request
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %s", err)
		}
	}()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// validate ensures the SOCKS5 proxy host and port are valid.
func (s *Socks5Client) validate(host, port string) error {
	if host == "" || port == "" {
		return errors.New("proxy host or port is empty")
	}

	// Validate port number
	if _, err := strconv.Atoi(port); err != nil {
		return errors.New("proxy port is invalid")
	}
	return nil
}
