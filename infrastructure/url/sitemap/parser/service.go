package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Service is responsible for parsing HTML content and extracting URLs.
type Service struct{}

// NewService creates and returns a new Parser service instance.
func NewService() *Service {
	return &Service{}
}

// URL represents the structure of each <url> element in the XML.
type URL struct {
	Location string `xml:"loc"`
}

// Sitemap represents the structure of the sitemap XML.
type Sitemap struct {
	URLs []URL `xml:"url"`
}

// Parse extracts URLs from the provided HTML content.
func (s *Service) Parse(body io.Reader) ([]string, error) {
	var sitemap Sitemap
	if err := xml.NewDecoder(body).Decode(&sitemap); err != nil {
		return nil, fmt.Errorf("parse sitemap: %w", err)
	}

	// If no URLs are found, return an error
	if len(sitemap.URLs) == 0 {
		return nil, errors.New("no URLs found in sitemap")
	}

	// Extract URLs from the sitemap
	urls := make([]string, 0, len(sitemap.URLs))
	for _, url := range sitemap.URLs {
		if s.isValid(url.Location) {
			urls = append(urls, url.Location)
		}
	}

	return urls, nil
}

// isValid filters URLs to include only job offer URLs relevant to Golang.
func (s *Service) isValid(url string) bool {
	// Check if the URL is a job offer and matches Golang-specific patterns
	return url != "" &&
		strings.Contains(url, "/job-offer/") &&
		(strings.Contains(url, "golang") || strings.Contains(url, "-go-"))
}

// SaveToFile saves the extracted URLs to a specified file for analysis.
func (s *Service) SaveToFile(urls []string, filePath string) error {
	// Open the file for writing, creating it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Printf("close file: %v", err)
		}
	}()

	// Write URLs to the file
	for _, url := range urls {
		if _, err = file.WriteString(url + "\n"); err != nil {
			return fmt.Errorf("write to file: %w", err)
		}
	}
	return nil
}
