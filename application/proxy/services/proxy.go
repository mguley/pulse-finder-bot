package services

import (
	"fmt"
	"infrastructure/http/client"
	"net/http"
	"sync"
	"time"
)

// Service provides functionality for managing HTTP clients configured to use a SOCKS5 proxy.
type Service struct {
	factory    *client.Factory // Factory for creating HTTP clients.
	host       string          // Hostname or IP address of the SOCKS5 proxy.
	port       string          // Port number of the SOCKS5 proxy.
	timeout    time.Duration   // Timeout duration for proxy interactions.
	mutex      sync.Mutex      // Mutex to ensure thread-safe access to the HTTP client.
	httpClient *http.Client    // Cached HTTP client instance.
}

// NewService creates a new Service instance.
func NewService(f *client.Factory, h, p string, t time.Duration) *Service {
	return &Service{factory: f, host: h, port: p, timeout: t, mutex: sync.Mutex{}, httpClient: nil}
}

// HttpClient retrieves or initializes an HTTP client configured to use the SOCKS5 proxy.
func (s *Service) HttpClient() (*http.Client, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// If HTTP client already exists, return it.
	if s.httpClient != nil {
		return s.httpClient, nil
	}

	// Create a new HTTP client using the factory.
	c, err := s.factory.CreateSocks5Client(s.host, s.port, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	// Cache the HTTP client for reuse.
	s.httpClient = c
	return s.httpClient, nil
}

// Close releases resources used by the HTTP client.
func (s *Service) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.httpClient != nil {
		// Close idle connections to release resources.
		t, ok := s.httpClient.Transport.(*http.Transport)
		if ok {
			t.CloseIdleConnections()
		}

		// Clear the cached HTTP client.
		s.httpClient = nil
	}
}
