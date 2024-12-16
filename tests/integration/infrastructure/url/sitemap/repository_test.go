package sitemap

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRepository_SaveUrls_Valid validates saving a list of valid URLs.
func TestRepository_SaveUrls_Valid(t *testing.T) {
	container := SetupTestContainer(t)
	sitemapRepo := container.SitemapRepository.Get()
	urlRepo := container.UrlRepository.Get()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	urls := []string{
		"https://example.com/job-offer/12345",
		"https://example.com/job-offer/67890",
	}

	// Save the URLs.
	err := sitemapRepo.SaveUrls(ctx, urls)
	require.NoError(t, err, "Repository should save valid URLs without errors")

	// Verify that the URLs are stored in the database.
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	list, err := urlRepo.FetchBatch(ctx, "pending", 2)
	require.NoError(t, err, "Repository should fetch URLs without errors")
	assert.Len(t, list, len(urls), "Number of saved URLs does not match input")
}

// TestRepository_SaveUrls_Empty validates saving an empty list of URLs.
func TestRepository_SaveUrls_Empty(t *testing.T) {
	container := SetupTestContainer(t)
	sitemapRepo := container.SitemapRepository.Get()
	urlRepo := container.UrlRepository.Get()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var urls []string

	// Save an empty list of URLs.
	err := sitemapRepo.SaveUrls(ctx, urls)
	require.NoError(t, err, "Repository should handle empty URL list without errors")

	// Verify that no new URLs are stored in the database.
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	list, err := urlRepo.FetchBatch(ctx, "pending", 2)
	require.NoError(t, err, "Repository should not return errors")
	assert.Len(t, list, len(urls), "Number of saved URLs does not match input")
}
