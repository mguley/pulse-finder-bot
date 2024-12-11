package client

import (
	"domain/useragent"
	"infrastructure/proxy/client"
	"net/http"
	"time"
)

// Factory manages the creation of HTTP clients.
type Factory struct {
	agent        useragent.Generator
	socks5Client *client.Socks5Client
}

// NewFactory creates a new instance of Factory.
func NewFactory(u useragent.Generator, c *client.Socks5Client) *Factory {
	return &Factory{agent: u, socks5Client: c}
}

// CreateSocks5Client returns an HTTP client configured to use a SOCKS5 proxy.
func (f *Factory) CreateSocks5Client(h, p string, timeout time.Duration) (*http.Client, error) {
	return f.socks5Client.HttpClient(h, p, timeout)
}

// CreateDefaultClient returns a standard HTTP client.
func (f *Factory) CreateDefaultClient() *http.Client {
	transport := &http.Transport{}
	uTransport := &RoundTripWithUserAgent{
		rt:        transport,
		userAgent: f.agent.Generate(),
	}

	return &http.Client{Transport: uTransport, Timeout: 10 * time.Second}
}
