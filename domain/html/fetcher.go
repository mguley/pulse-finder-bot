package html

import "context"

// Fetcher defines the contract for fetching HTML content.
type Fetcher interface {
	// Fetch retrieves the raw HTML content from the specified URL.
	// Returns the raw HTML content, or an error if the fetching process fails.
	Fetch(ctx context.Context, url string) (body string, err error)
}
