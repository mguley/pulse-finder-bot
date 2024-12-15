package html

import "context"

// Fetcher defines the interface for fetching raw HTML content from a given URL.
type Fetcher interface {
	// Fetch retrieves the raw HTML content from the specified URL.
	// Returns the raw HTML content, or an error if the fetching process fails.
	Fetch(ctx context.Context, url string) (string, error)

	// Close releases resources used by the HTTP client.
	Close()
}
