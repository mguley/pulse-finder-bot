package source

import "context"

// Handler defines the operations required for managing a source.
type Handler interface {
	// ProcessURLs retrieves a list of URLs from the source and saves them into a data source.
	// Returns an error if the operation fails.
	ProcessURLs(ctx context.Context) error

	// ProcessHTML processes the HTML content of URL and extracts vacancy details.
	// Returns an error if the operation fails.
	ProcessHTML(ctx context.Context, batchSize int) error
}
