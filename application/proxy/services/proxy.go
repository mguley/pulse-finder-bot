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
func NewService(factory *client.Factory, host, port string, timeout time.Duration) *Service {
	return &Service{factory: factory, host: host, port: port, timeout: timeout, mutex: sync.Mutex{}, httpClient: nil}
}

// HttpClient provides HTTP client configured to use the SOCKS5 proxy.
func (s *Service) HttpClient() (client *http.Client, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.httpClient != nil {
		return s.httpClient, nil
	}

	httpClient, err := s.factory.CreateSocks5Client(s.host, s.port, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("create proxy http client: %v", err)
	}

	// Cache the HTTP client for reuse.
	s.httpClient = httpClient
	return s.httpClient, nil
}

// Close releases resources.
func (s *Service) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.httpClient != nil {
		t, ok := s.httpClient.Transport.(*http.Transport)
		if ok {
			t.CloseIdleConnections()
		}
		s.httpClient = nil
	}
}
