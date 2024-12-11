package agent

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// ChromeUserAgentGenerator generates User-Agent strings simulating Google Chrome browsers.
type ChromeUserAgentGenerator struct {
	versions         []string
	operatingSystems []string
}

// NewChromeUserAgentGenerator creates a new instance of ChromeUserAgentGenerator.
func NewChromeUserAgentGenerator() *ChromeUserAgentGenerator {
	return &ChromeUserAgentGenerator{
		versions: []string{
			"126.0.6478.114", "126.0.6478.62", "126.0.6478.61",
			"126.0.6478.56", "124.0.6367.243", "124.0.6367.233",
			"124.0.6367.230", "124.0.6367.221", "124.0.6367.208",
			"124.0.6367.201", "124.0.6367.118", "123.0.6358.132",
			"123.0.6358.121", "122.0.6345.98", "122.0.6345.67",
		},
		operatingSystems: []string{
			"Windows NT 10.0; Win64; x64",
			"Macintosh; Intel Mac OS X 10_15_7",
			"X11; Linux x86_64", "Windows NT 6.1; Win64; x64",
			"Macintosh; Intel Mac OS X 10_14_6",
		},
	}
}

// Generate creates a random User-Agent string for Chrome browsers.
func (g *ChromeUserAgentGenerator) Generate() string {
	v := g.versions[g.rand(len(g.versions))]
	os := g.operatingSystems[g.rand(len(g.operatingSystems))]

	return fmt.Sprintf("Mοzilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, v)
}

// rand generates random integer.
func (g *ChromeUserAgentGenerator) rand(m int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(m)))
	if err != nil {
		return m
	}
	return int(n.Int64())
}
