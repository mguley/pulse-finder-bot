package alfa

import (
	"application/dependency"
	"domain/html"
	htmlAlfa "infrastructure/html/source/alfa"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	AlfaHtmlParser dependency.LazyDependency[html.Parser]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.AlfaHtmlParser = dependency.LazyDependency[html.Parser]{
		InitFunc: func() html.Parser {
			return htmlAlfa.NewParser()
		},
	}

	return c
}
