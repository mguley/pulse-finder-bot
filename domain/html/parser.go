package html

import "application/url/processor/dto"

// Parser defines the interface for parsing raw HTML content into structured data.
type Parser interface {
	// Parse processes the raw HTML content and converts it into a Vacancy DTO.
	// Returns a DTO populated with the parsed data, or an error if the parsing process fails.
	Parse(html string) (*dto.Vacancy, error)
}
