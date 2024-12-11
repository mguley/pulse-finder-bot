package useragent

// Generator defines the interface for generating User-Agent strings.
type Generator interface {
	// Generate returns a User-Agent string that simulates a browser and device.
	Generate() string
}
