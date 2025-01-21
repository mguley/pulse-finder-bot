package alfa

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// HttpClient defines an interface for making HTTP requests and managing connections.
type HttpClient interface {
	// Do sends an HTTP request and returns an HTTP response.
	Do(req *http.Request) (response *http.Response, err error)

	// CloseIdleConnections closes any idle connections in a "keep-alive" state.
	CloseIdleConnections()
}

// Fetcher fetches data from URL using an HTTP client.
type Fetcher struct {
	httpClient  HttpClient // httpClient is a client used to send requests.
	maxBodySize int64      // maxBodySize is a number of bytes to read from the response body.
}

// NewFetcher creates a new Fetcher instance with a configurable body size limit.
func NewFetcher(httpClient HttpClient, maxBodySize int64) *Fetcher {
	return &Fetcher{httpClient: httpClient, maxBodySize: maxBodySize}
}

// Fetch sends an HTTP GET request to the specified URL and returns the response body as a string.
func (f *Fetcher) Fetch(ctx context.Context, url string) (result string, err error) {
	var (
		request  *http.Request
		response *http.Response
		body     []byte
	)

	if request, err = http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody); err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	if response, err = f.httpClient.Do(request); err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			fmt.Printf("close response body: %v", err)
		}
	}()
	defer f.httpClient.CloseIdleConnections()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status code: %d", response.StatusCode)
	}

	limitedReader := io.LimitReader(response.Body, f.maxBodySize)
	if body, err = io.ReadAll(limitedReader); err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	return string(body), nil
}
