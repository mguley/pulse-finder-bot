package sitemap

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_Parse_ValidHTML validates that the parser extracts URLs correctly from valid HTML.
func TestParser_Parse_ValidHTML(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.SitemapParser.Get()

	// Sample valid HTML content with footer and job offer links.
	htmlContent := `
		<html>
			<footer>
				<a href="https://example.com">Host</a>
			</footer>
			<div class="MuiBox-root">
				<a href="/job-offer/12345">Job 1</a>
				<a href="/job-offer/67890">Job 2</a>
			</div>
		</html>
	`

	body := strings.NewReader(htmlContent)

	// Parse the content
	urls, err := parser.Parse(body)
	require.NoError(t, err, "Parser failed to process valid HTML")

	// Assert the extracted URLs
	expectedUrls := []string{
		"https://example.com/job-offer/12345",
		"https://example.com/job-offer/67890",
	}
	assert.ElementsMatch(t, expectedUrls, urls, "Parsed URLs do not match expected values")
}

// TestParser_Parse_MissingHost validates that the parser handles missing host gracefully.
func TestParser_Parse_MissingHost(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.SitemapParser.Get()

	// HTML content without a footer or host link.
	htmlContent := `
		<html>
			<div class="MuiBox-root">
				<a href="/job-offer/12345">Job 1</a>
				<a href="/job-offer/67890">Job 2</a>
			</div>
		</html>
	`

	body := strings.NewReader(htmlContent)

	// Parse the content and expect an error.
	urls, err := parser.Parse(body)
	require.Error(t, err, "Parser should fail when host is missing")
	assert.Nil(t, urls, "URLs should be nil when parsing fails")
	assert.Contains(t, err.Error(), "failed to find host", "Error message does not indicate missing host")
}

// TestParser_Parse_EmptyHTML validates that the parser handles empty HTML input.
func TestParser_Parse_EmptyHTML(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.SitemapParser.Get()

	// Empty HTML content.
	htmlContent := ""
	body := strings.NewReader(htmlContent)

	// Parse the content and expect an error.
	urls, err := parser.Parse(body)
	require.Error(t, err, "Parser should fail on empty HTML")
	assert.Nil(t, urls, "URLs should be nil when parsing fails")
}
