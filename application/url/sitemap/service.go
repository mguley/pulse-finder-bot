package sitemap

import (
	"context"
	"fmt"
	"infrastructure/url/sitemap/fetcher"
	"infrastructure/url/sitemap/notifier"
	"infrastructure/url/sitemap/parser"
	"infrastructure/url/sitemap/repository"
	"net/http"
)

// Service orchestrates the processing of URLs from a sitemap.
// It coordinates fetching, parsing, notifying and storing URLs.
type Service struct {
	fetcher  *fetcher.Service             // Service responsible for fetching HTML content.
	parser   *parser.Service              // Service for parsing HTML content and extracting URLs.
	repo     *repository.Service          // Service for storing extracted URLs into the data source.
	notifier *notifier.Service            // Service for handling notifications (e.g., logging proxy IPs).
	client   func() (*http.Client, error) // Function to provide an HTTP client.
}

// NewService creates and returns a new instance of the Sitemap service.
func NewService(
	f *fetcher.Service,
	p *parser.Service,
	repo *repository.Service,
	n *notifier.Service,
	client func() (*http.Client, error),
) *Service {
	return &Service{
		fetcher:  f,
		parser:   p,
		repo:     repo,
		notifier: n,
		client:   client,
	}
}

// ProcessUrls orchestrates the complete flow of fetching, parsing, notifying, and saving URLs.
func (s *Service) ProcessUrls(ctx context.Context, url string) error {
	// Fetch content from the URL.
	body, err := s.fetcher.Fetch(ctx, url)
	if err != nil {
		return fmt.Errorf("fetch url: %w", err)
	}
	defer func() {
		if err = body.Close(); err != nil {
			fmt.Printf("close body err:%v", err)
		}
	}()

	// Notify (log the proxy's IP address).
	client, err := s.client()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}
	if err = s.notifier.Notify(client); err != nil {
		return fmt.Errorf("notify: %w", err)
	}

	// Parse the fetched content to extract URLs.
	urls, err := s.parser.Parse(body)
	if err != nil {
		return fmt.Errorf("parse urls: %w", err)
	}

	// Save the extracted URLs to the data source.
	if err = s.repo.SaveUrls(ctx, urls); err != nil {
		return fmt.Errorf("save urls: %w", err)
	}

	return nil
}
