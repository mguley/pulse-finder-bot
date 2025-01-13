package alfa

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_ValidHTML tests parsing a valid HTML document.
func TestParser_ValidHTML(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.AlfaHtmlParser.Get()

	// Example valid HTML content
	htmlContent := `
		<html>
			<head><title>Backend Developer</title></head>
			<body>
				<div class="row align-items-center gx-1">
					<a class="text-reset">Tech Innovations</a>
				</div>
			</body>
		</html>
	`

	vacancy, err := parser.Parse(htmlContent)
	defer vacancy.Release()

	require.NoError(t, err, "Parser should not return an error for valid HTML")
	require.NotNil(t, vacancy, "Vacancy should not be nil for valid HTML")
	assert.Equal(t, "Backend Developer", vacancy.Title, "Vacancy title should match")
	assert.Equal(t, "Tech Innovations", vacancy.Company, "Vacancy company should match")
}

// TestParser_MissingElements tests parsing HTML with missing elements.
func TestParser_MissingElements(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.AlfaHtmlParser.Get()

	// HTML content missing the company and description
	htmlContent := `
		<html>
			<head><title>Frontend Engineer</title></head>
			<body></body>
		</html>
	`

	vacancy, err := parser.Parse(htmlContent)
	defer vacancy.Release()

	require.NoError(t, err, "Parser should not return an error for HTML with missing elements")
	require.NotNil(t, vacancy, "Vacancy should not be nil for HTML with missing elements")
	assert.Equal(t, "Frontend Engineer", vacancy.Title, "Vacancy title should match")
	assert.Equal(t, "Unknown Company", vacancy.Company, "Vacancy company should default to 'Unknown Company'")
}

// TestParser_MalformedHTML tests parsing a malformed HTML document.
func TestParser_MalformedHTML(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.AlfaHtmlParser.Get()

	// Malformed HTML content
	htmlContent := `
		<html>
			<head><title>DevOps Specialist</title></head>
			<body>
				<div class="row align-items-center gx-1
	`

	vacancy, err := parser.Parse(htmlContent)
	defer vacancy.Release()

	require.NoError(t, err, "Parser should not return an error for malformed HTML")
	require.NotNil(t, vacancy, "Vacancy should not be nil for malformed HTML")
	assert.Equal(t, "DevOps Specialist", vacancy.Title, "Vacancy title should match")
	assert.Equal(t, "Unknown Company", vacancy.Company, "Vacancy company should default to 'Unknown Company'")
}

// TestParser_EmptyHTML tests parsing an empty HTML document.
func TestParser_EmptyHTML(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.AlfaHtmlParser.Get()

	// Empty HTML content
	htmlContent := ""
	vacancy, err := parser.Parse(htmlContent)
	defer vacancy.Release()

	require.NoError(t, err, "Parser should not return an error for empty HTML")
	require.NotNil(t, vacancy, "Vacancy should not be nil for empty HTML")
	assert.Equal(t, "Unknown Title", vacancy.Title, "Vacancy title should default to 'Unknown Title'")
	assert.Equal(t, "Unknown Company", vacancy.Company, "Vacancy company should default to 'Unknown Company'")
}
