package beta

import (
	"application/proxy/services"
	"context"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient defines the interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Fetcher is responsible for fetching HTML content over HTTP using a proxy service.
type Fetcher struct {
	client      HTTPClient        // HTTP client used to make requests.
	maxBodySize int64             // Maximum size of the response body that can be read, in bytes.
	proxy       *services.Service // Proxy service to manage HTTP client lifecycle.
}

// NewFetcher creates a new Fetcher instance with a configurable body size limit.
func NewFetcher(proxy *services.Service, maxBodySize int64) (*Fetcher, error) {
	c, err := proxy.HttpClient()
	if err != nil {
		return nil, fmt.Errorf("create fetcher: %w", err)
	}
	return &Fetcher{client: c, maxBodySize: maxBodySize, proxy: proxy}, nil
}

// Fetch retrieves the HTML content from the specified URL.
func (f *Fetcher) Fetch(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Printf("close response body error: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status: %d", resp.StatusCode)
	}

	limitedReader := io.LimitReader(resp.Body, f.maxBodySize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}
	return string(body), nil
}

// Close releases resources used by the Fetcher.
func (f *Fetcher) Close() {
	if f.proxy != nil {
		f.proxy.Close()
	}
}
