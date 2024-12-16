package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Service handles fetching HTML content from a given URL.
type Service struct {
	clientProvider func() (*http.Client, error) // Provides an HTTP client.
}

// NewService creates and returns a new instance of Fetcher service.
func NewService(clientProvider func() (*http.Client, error)) *Service {
	return &Service{clientProvider: clientProvider}
}

// Fetch retrieves the content of the given URL.
func (s *Service) Fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	client, err := s.clientProvider()
	if err != nil {
		return nil, fmt.Errorf("http client: %w", err)
	}

	// Create the HTTP request.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Perform the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
