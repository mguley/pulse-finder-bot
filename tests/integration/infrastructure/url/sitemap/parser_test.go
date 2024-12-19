package sitemap

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_Parse_ValidXML validates that the parser extracts URLs correctly from valid XML content.
func TestParser_Parse_ValidXML(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.SitemapParser.Get()

	// Sample valid XML content with job offer links.
	xmlContent := `
		<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
			<url>
				<loc>https://example.com/job-offer/12-go-12345</loc>
			</url>
			<url>
				<loc>https://example.com/job-offer/12-go-67890</loc>
			</url>
		</urlset>
	`

	body := strings.NewReader(xmlContent)

	// Parse the content
	urls, err := parser.Parse(body)
	require.NoError(t, err, "Parser failed to process valid XML")

	// Assert the extracted URLs
	expectedUrls := []string{
		"https://example.com/job-offer/12-go-12345",
		"https://example.com/job-offer/12-go-67890",
	}
	assert.ElementsMatch(t, expectedUrls, urls, "Parsed URLs do not match expected values")
}

// TestParser_Parse_NoJobOffers validates that the parser handles XML without valid job offer URLs.
func TestParser_Parse_NoJobOffers(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.SitemapParser.Get()

	// XML content without any job offer links.
	xmlContent := `
		<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
			<url>
				<loc>https://example.com/other-page</loc>
			</url>
		</urlset>
	`

	body := strings.NewReader(xmlContent)

	// Parse the content and expect no valid URLs.
	urls, err := parser.Parse(body)
	require.NoError(t, err, "Parser should not fail when no job offers exist")
	assert.Empty(t, urls, "URLs should be empty when no job offer links exist")
}

// TestParser_Parse_EmptyXML validates that the parser handles empty XML content.
func TestParser_Parse_EmptyXML(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.SitemapParser.Get()

	// Empty XML content.
	xmlContent := ""
	body := strings.NewReader(xmlContent)

	// Parse the content and expect an error.
	urls, err := parser.Parse(body)
	require.Error(t, err, "Parser should fail on empty XML")
	assert.Nil(t, urls, "URLs should be nil when parsing fails")
	assert.Contains(t, err.Error(), "parse sitemap", "Error message does not indicate XML parsing failure")
}
