package beta

import "testing"

// SetupTestContainer initializes the TestContainer and handles cleanup.
func SetupTestContainer(t *testing.T) *TestContainer {
	c := NewTestContainer()

	t.Cleanup(func() {
		fetcher := c.BetaHtmlFetcher.Get()
		fetcher.Close()
	})

	return c
}
