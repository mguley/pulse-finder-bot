package beta

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
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses the provided raw HTML content and extracts vacancy details.
func (p *Parser) Parse(htmlContent string) (*dto.Vacancy, error) {
	// Load the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to load HTML document: %w", err)
	}

	// Extract the title from the <title> tag
	title := strings.TrimSpace(doc.Find("title").Text())
	if title == "" {
		title = "Unknown Title"
	}

	// Extract the company name from the specified selector
	company := strings.TrimSpace(doc.Find("h2.MuiTypography-root").First().Text())
	if company == "" {
		company = "Unknown Company"
	}

	// Extract the job description
	_ = extractDescription(doc)
	description := "-" // todo: improve

	location := extractLocation(doc)

	// Populate the vacancy DTO with extracted data
	v := dto.GetVacancy()
	v.Title = title
	v.Company = company
	v.Description = description
	v.Location = location
	v.PostedAt = time.Now()

	return v, nil
}

// extractLocation extracts the location from the HTML document.
func extractLocation(doc *goquery.Document) string {
	var location string

	// Find all <div> elements with the "MuiBox-root" class
	doc.Find("div.MuiBox-root").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// Get the child <div> elements
		childDivs := s.ChildrenFiltered("div")

		// Ensure the <div> has exactly two child <div> elements
		if childDivs.Length() == 2 {
			// Check if the first child <div> has the text "Operating mode"
			firstChild := childDivs.Eq(0)
			if strings.TrimSpace(firstChild.Text()) != "Operating mode" {
				return true // Continue if the text is not "Operating mode"
			}

			// Extract the text from the second child <div>
			secondChild := childDivs.Eq(1)
			location = strings.TrimSpace(secondChild.Text())
			return false // Break the loop as we've found the desired element
		}
		return true // Continue loop if conditions are not met
	})

	if location == "" {
		location = "Unknown Location" // Default value if not found
	}
	return location
}

// extractDescription extracts the job description from the HTML document.
func extractDescription(doc *goquery.Document) string {
	var description string

	// Narrow down to <div> elements with a specific class ("MuiBox-root")
	doc.Find("div.MuiBox-root").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// Get the child <div> elements
		childDivs := s.ChildrenFiltered("div")

		// Ensure the <div> has exactly two child <div> elements
		if childDivs.Length() == 2 {
			// Check if the first child <div> has an <h3> as the first child
			firstChild := childDivs.Eq(0)
			if firstChild.ChildrenFiltered("h3").Length() == 0 {
				return true // Continue if no <h3> in the first child <div>
			}

			// Check if the second child <div> has a <p> as the first child
			secondChild := childDivs.Eq(1)
			if secondChild.ChildrenFiltered("p").Length() == 0 {
				return true // Continue if no <p> in the second child <div>
			}

			// Extract HTML content of the second child <div>
			var err error
			description, err = secondChild.Html()
			if err != nil {
				description = ""
			}
			return false // Break the loop as we've found the desired element
		}
		return true // Continue loop if conditions are not met
	})

	return description
}
