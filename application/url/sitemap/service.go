package sitemap

import (
	"context"
	"fmt"
	"infrastructure/url/sitemap/fetcher"
	"infrastructure/url/sitemap/notifier"
	"infrastructure/url/sitemap/repository"
	"io"
	"net/http"
)

// Parser defines the contract for parser.
type Parser interface {
	Parse(body io.Reader) ([]string, error)
}

// Option defines a functional option for configuring the Sitemap Service.
type Option func(service *Service)

// Service orchestrates the processing of URLs from a sitemap.
// It coordinates fetching, parsing, notifying and storing URLs.
type Service struct {
	fetcher  *fetcher.Service             // Service responsible for fetching HTML content.
	parser   Parser                       // Service for parsing HTML content and extracting URLs.
	repo     *repository.Service          // Service for storing extracted URLs into the data source.
	notifier *notifier.Service            // Service for handling notifications (e.g., logging proxy IPs).
	client   func() (*http.Client, error) // Function to provide an HTTP client.
}

// NewService creates and returns a new instance of the Sitemap service.
func NewService(options ...Option) *Service {
	s := &Service{}
	for _, option := range options {
		option(s)
	}
	return s
}

// WithFetcher sets the fetcher dependency.
func WithFetcher(f *fetcher.Service) Option {
	return func(s *Service) {
		s.fetcher = f
	}
}

// WithNotifier sets the notifier dependency.
func WithNotifier(n *notifier.Service) Option {
	return func(s *Service) {
		s.notifier = n
	}
}

// WithRepository sets the repository dependency.
func WithRepository(r *repository.Service) Option {
	return func(s *Service) {
		s.repo = r
	}
}

// WithHTTPClient sets the HTTP client provider.
func WithHTTPClient(client func() (*http.Client, error)) Option {
	return func(s *Service) {
		s.client = client
	}
}

// WithParser sets the parser dependency (XML or RSS).
func WithParser(parser Parser) Option {
	return func(s *Service) {
		s.parser = parser
	}
}

// ProcessUrls orchestrates the complete flow of fetching, parsing, notifying, and saving URLs.
func (s *Service) ProcessUrls(ctx context.Context, url string) error {
	// Notify (log the proxy's IP address).
	client, err := s.client()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}
	if err = s.notifier.Notify(ctx, client); err != nil {
		return fmt.Errorf("notify: %w", err)
	}

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
