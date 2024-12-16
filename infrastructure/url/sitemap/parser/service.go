package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Service is responsible for parsing HTML content and extracting URLs.
type Service struct{}

// NewService creates and returns a new Parser service instance.
func NewService() *Service {
	return &Service{}
}

// Parse extracts URLs from the provided HTML content.
func (s *Service) Parse(body io.Reader) ([]string, error) {
	// Read the entire body into a string.
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	htmlContent := string(bodyBytes)

	// Load the HTML document using goquery.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to load HTML document: %w", err)
	}

	var urls []string
	var host string

	// Extract the host from the footer (or another location).
	doc.Find("footer a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if exists && strings.HasPrefix(href, "https://") {
			host = href
			return false // Break after finding the first valid host.
		}
		return true
	})

	// If no host is found, return an error.
	if host == "" {
		return nil, fmt.Errorf("failed to find host in HTML document")
	}

	// Extract and combine job offer URLs with the host.
	doc.Find(".MuiBox-root a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, "/job-offer/") { // Filter job offer URLs.
			fullUrl := host + href
			urls = append(urls, fullUrl)
		}
	})

	return urls, nil
}
