package beta

import (
	"application/url/sitemap"
	"context"
	"fmt"
)

// Handler processes URLs and HTML content.
type Handler struct {
	url            string           // Base URL of the sitemap to process.
	sitemapService *sitemap.Service // Service orchestrates the processing of URLs from a sitemap.
}

// NewHandler creates and returns a new Handler instance.
func NewHandler(url string, s *sitemap.Service) *Handler {
	return &Handler{url: url, sitemapService: s}
}

// ProcessURLs retrieves and processes sitemap URLs.
func (h *Handler) ProcessURLs(ctx context.Context) error {
	if err := h.sitemapService.ProcessUrls(ctx, h.url); err != nil {
		return fmt.Errorf("process urls: %w", err)
	}
	return nil
}

// ProcessHTML processes URLs in batches with a delay.
func (h *Handler) ProcessHTML(ctx context.Context, batchSize int) error {
	// todo
	return nil
}
