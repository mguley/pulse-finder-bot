package sitemap

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFetcher_Fetch_ValidURL validates fetching content from a valid URL.
func TestFetcher_Fetch_ValidURL(t *testing.T) {
	container := SetupTestContainer(t)
	fetcher := container.SitemapFetcher.Get()
	url := container.Config.Get().Proxy.PingUrl

	// Fetch the content
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := fetcher.Fetch(ctx, url)

	require.NoError(t, err, "Fetcher should not return an error for a valid URL")
	defer func() {
		if err = body.Close(); err != nil {
			require.NoError(t, err, "Should not error")
		}
	}()

	// Read the response body.
	content, err := io.ReadAll(body)
	require.NoError(t, err, "Failed to read response body")
	assert.NotEmpty(t, content, "Response body should not be empty")
}

// TestFetcher_Fetch_InvalidURL validates handling of an invalid URL.
func TestFetcher_Fetch_InvalidURL(t *testing.T) {
	container := SetupTestContainer(t)
	fetcher := container.SitemapFetcher.Get()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define an invalid URL.
	url := "http://invalid.invalid-url"

	body, err := fetcher.Fetch(ctx, url)

	// Expect an error to occur.
	require.Error(t, err, "Expected fetcher to fail with invalid URL")
	assert.Nil(t, body, "Expected no body to be returned for invalid URL")
}
