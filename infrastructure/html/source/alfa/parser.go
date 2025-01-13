package alfa

import (
	"application/url/processor/dto"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Parser is responsible for extracting vacancy details from raw HTML content.
type Parser struct{}

// NewParser creates and returns a new Parser instance.
func NewParser() *Parser { return &Parser{} }

// Parse parses the provided raw HTML content and extracts vacancy details.
func (p *Parser) Parse(htmlContent string) (*dto.Vacancy, error) {
	// Load the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("load HTML document: %w", err)
	}

	// Extract the title from the <title> tag
	title := strings.TrimSpace(doc.Find("title").Text())
	if title == "" {
		title = "Unknown Title"
	}

	// Extract the company name from the specified selector
	company := strings.TrimSpace(doc.Find("div.row.align-items-center.gx-1 a.text-reset").Text())
	if company == "" {
		company = "Unknown Company"
	}

	// Extract the job description
	description := "-" // todo: improve
	location := "-"

	// Populate the vacancy DTO with extracted data
	v := dto.GetVacancy()
	v.Title = title
	v.Company = company
	v.Description = description
	v.Location = location
	v.PostedAt = time.Now()

	return v, nil
}
