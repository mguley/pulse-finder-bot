package beta

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_ValidHTML tests parsing a valid HTML document.
func TestParser_ValidHTML(t *testing.T) {
	t.Skip("Skipping TestParser_ValidHTML temporarily, parser logic requires improvements (description field)")

	container := SetupTestContainer(t)
	parser := container.BetaHtmlParser.Get()

	// Example valid HTML content
	htmlContent := `
		<html>
			<head><title>Software Engineer</title></head>
			<body>
				<p class="MuiTypography-root MuiTypography-h3">Tech Corp</p>
				<div class="MuiBox-root">
					<div>
						<h3>Job Description</h3>
					</div>
					<div>
						<p>We are looking for a Software Engineer to join our team.</p>
					</div>
				</div>
			</body>
		</html>
	`

	vacancy, err := parser.Parse(htmlContent)
	defer vacancy.Release()

	require.NoError(t, err, "Parser should not return an error for valid HTML")
	require.NotNil(t, vacancy, "Vacancy should not be nil for valid HTML")
	assert.Equal(t, "Software Engineer", vacancy.Title, "Vacancy title should match")
	assert.Equal(t, "Tech Corp", vacancy.Company, "Vacancy description should match")
	assert.Contains(t, vacancy.Description, "We are looking for a Software Engineer", "Vacancy description should match")
}

// TestParser_MissingElements tests parsing HTML with missing elements.
func TestParser_MissingElements(t *testing.T) {
	t.Skip("Skipping TestParser_MissingElements temporarily, parser logic requires improvements (description field)")

	container := SetupTestContainer(t)
	parser := container.BetaHtmlParser.Get()

	// HTML content missing company and description
	htmlContent := `
		<html>
			<head><title>Product Manager</title></head>
			<body></body>
		</html>
	`

	vacancy, err := parser.Parse(htmlContent)
	defer vacancy.Release()

	require.NoError(t, err, "Parser should not return an error for HTML with missing elements")
	require.NotNil(t, vacancy, "Vacancy should not be nil for HTML with missing elements")
	assert.Equal(t, "Product Manager", vacancy.Title, "Vacancy title should match")
	assert.Empty(t, vacancy.Company, "Vacancy company should be empty")
	assert.Empty(t, vacancy.Description, "Vacancy description should be empty")
}

// TestParser_MalformedHTML tests parsing a malformed HTML document.
func TestParser_MalformedHTML(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.BetaHtmlParser.Get()

	// Malformed HTML content
	htmlContent := `
		<html>
			<head><title>Data Scientist</title></head>
			<body>
				<p class="MuiTypography-root
	`

	vacancy, err := parser.Parse(htmlContent)
	defer vacancy.Release()

	require.NoError(t, err, "Parser should not return an error for malformed HTML")
	require.NotNil(t, vacancy, "Vacancy should not be nil for malformed HTML")
	assert.Equal(t, vacancy.Title, "Data Scientist", "Vacancy title should match")
}

// TestParser_EmptyHTML tests parsing an empty HTML document.
func TestParser_EmptyHTML(t *testing.T) {
	container := SetupTestContainer(t)
	parser := container.BetaHtmlParser.Get()

	// Empty HTML content
	htmlContent := ""
	vacancy, err := parser.Parse(htmlContent)
	defer vacancy.Release()

	require.NoError(t, err, "Parser should not return an error for empty HTML")
	require.NotNil(t, vacancy, "Vacancy should not be nil for empty HTML")
}
