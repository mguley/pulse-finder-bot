package circuit

import (
	"application/proxy/services"
	"context"
	"fmt"
	"infrastructure/proxy"
	"net/http"
)

// Manager handles circuit changes and verification by interacting with identity and proxy services.
type Manager struct {
	identityService *services.Identity // Identity service for requesting a new circuit.
	proxyService    *services.Service  // Proxy service for routing traffic through the proxy.
	pingUrl         string             // URL used to verify the new circuit.
}

// NewManager creates a new Manager instance.
func NewManager(i *services.Identity, p *services.Service, u string) *Manager {
	return &Manager{identityService: i, proxyService: p, pingUrl: u}
}

// ChangeCircuit requests a new circuit via the identity service and verifies it.
func (m *Manager) ChangeCircuit() (string, error) {
	// Request a new circuit from the identity service.
	if err := m.identityService.Request(); err != nil {
		return "", fmt.Errorf("identity request: %w", err)
	}

	// Retrieve an HTTP client from the proxy service.
	client, err := m.proxyService.HttpClient()
	if err != nil {
		return "", fmt.Errorf("http client: %w", err)
	}
	defer m.proxyService.Close()

	// Verify the circuit change using the provided HTTP client.
	ip, err := m.verify(client)
	if err != nil {
		return "", fmt.Errorf("verify: %w", err)
	}
	return ip, nil
}

// verify validates the new circuit by performing an HTTP GET request to the ping URL.
func (m *Manager) verify(client *http.Client) (string, error) {
	checker := proxy.GetChecker().Reset()
	defer checker.Release()

	checker.SetClient(client)

	// Perform the verification request.
	ip, err := checker.GetInfo(context.Background(), m.pingUrl)
	if err != nil {
		return "", fmt.Errorf("get info: %w", err)
	}
	return ip, nil
}
