package circuit

import (
	"application/proxy/commands"
	"application/proxy/services"
	"context"
	"fmt"
	"net/http"
	"sync"
)

// Manager handles circuit changes and verification by interacting with identity and proxy services.
type Manager struct {
	identityService *services.Identity // identityService manages requests to change proxy circuit.
	proxyService    *services.Service  // proxyService manages HTTP clients configured to use SOCKS5 proxy.
	mutex           sync.Mutex         // mutex to ensure thread-save access.
	url             string             // url is the target endpoint for the HTTP GET request.
}

// NewManager creates a new Manager instance.
func NewManager(identity *services.Identity, proxy *services.Service, url string) *Manager {
	return &Manager{identityService: identity, proxyService: proxy, url: url}
}

// ChangeCircuit requests a new circuit and validates it.
func (m *Manager) ChangeCircuit() (result string, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var (
		client   *http.Client
		response string
	)

	if err = m.identityService.Change(); err != nil {
		return "", fmt.Errorf("identity change: %w", err)
	}

	if client, err = m.proxyService.HttpClient(); err != nil {
		return "", fmt.Errorf("get http client: %w", err)
	}
	defer m.proxyService.Close()

	if response, err = m.validate(client); err != nil {
		return "", fmt.Errorf("validate: %w", err)
	}
	return response, nil
}

// validate validates a new circuit by performing HTTP request to the validation URL.
func (m *Manager) validate(client *http.Client) (result string, err error) {
	check := commands.GetCheckCommand()
	defer check.Release()

	check.SetClient(client, m.url)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if result, err = check.Execute(ctx); err != nil {
		return "", fmt.Errorf("execute check: %w", err)
	}
	return result, nil
}
