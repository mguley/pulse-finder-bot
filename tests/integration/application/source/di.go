package source

import (
	"application/dependency"
	"application/source"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	SourceFactory dependency.LazyDependency[*source.Factory]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.SourceFactory = dependency.LazyDependency[*source.Factory]{
		InitFunc: source.NewFactory,
	}

	return c
}
