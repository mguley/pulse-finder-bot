package sitemap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNotifier_Notify validates the notifier's behavior with a valid proxy-configured HTTP client.
func TestNotifier_Notify(t *testing.T) {
	container := SetupTestContainer(t)
	notifier := container.SitemapNotifier.Get()
	proxyService := container.ProxyService.Get()

	// Obtain the proxy configured HTTP client.
	client, err := proxyService.HttpClient()
	require.NoError(t, err, "Failed to create proxy-configured HTTP client")

	// Call Notify to log the proxy's IP address.
	err = notifier.Notify(client)
	require.NoError(t, err, "Notifier should not return an error for a valid client")
}
