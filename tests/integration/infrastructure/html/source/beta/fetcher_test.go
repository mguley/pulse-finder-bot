package beta

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFetcher_FetchValidURL tests fetching a valid URL.
func TestFetcher_FetchValidURL(t *testing.T) {
	container := SetupTestContainer(t)
	fetcher := container.BetaHtmlFetcher.Get()
	url := container.Config.Get().Proxy.PingUrl

	// Fetch the content
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	content, err := fetcher.Fetch(ctx, url)

	require.NoError(t, err, "Fetcher should not return an error for a valid URL")
	assert.Contains(t, content, "ip", "Response should contain the key 'ip'")
}

// TestFetcher_FetchInvalidURL tests fetching an invalid URL.
func TestFetcher_FetchInvalidURL(t *testing.T) {
	container := SetupTestContainer(t)
	fetcher := container.BetaHtmlFetcher.Get()

	// Invalid URL to fetch
	url := "http://invalid-url-test"

	// Fetch the content
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	content, err := fetcher.Fetch(ctx, url)

	require.Error(t, err, "Fetcher should return an error for an invalid URL")
	assert.Empty(t, content, "Content should be empty for an invalid URL")
}
