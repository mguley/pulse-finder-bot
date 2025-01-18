package client

import (
	"domain/useragent"
	"infrastructure/proxy/client"
	"net/http"
	"time"
)

// Factory manages the creation of HTTP clients.
type Factory struct {
	agent           useragent.Generator
	socks5Client    *client.Socks5Client
	idleConnTimeout time.Duration
}

// NewFactory creates a new instance of Factory.
func NewFactory(agent useragent.Generator, client *client.Socks5Client) *Factory {
	return &Factory{agent: agent, socks5Client: client, idleConnTimeout: time.Duration(10) * time.Second}
}

// CreateSocks5Client returns an HTTP client configured to use a SOCKS5 proxy.
func (f *Factory) CreateSocks5Client(host, port string, timeout time.Duration) (client *http.Client, err error) {
	return f.socks5Client.HttpClient(host, port, timeout)
}

// CreateDefaultClient returns a default HTTP client.
func (f *Factory) CreateDefaultClient(timeout time.Duration) *http.Client {
	transport := &http.Transport{
		IdleConnTimeout: f.idleConnTimeout,
	}
	uTransport := &RoundTripWithUserAgent{
		rt:        transport,
		userAgent: f.agent.Generate(),
	}
	return &http.Client{Transport: uTransport, Timeout: timeout}
}
