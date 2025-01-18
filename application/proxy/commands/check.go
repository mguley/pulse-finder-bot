package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// checkCommandPool provided on-demand access to a sync.Pool.
var checkCommandPool = commandPool()

// commandPool returns a function that returns a *sync.Pool.
// The returned function uses sync.Once to ensure the *sync.Pool is created exactly once.
func commandPool() func() *sync.Pool {
	var (
		once sync.Once
		pool *sync.Pool
	)
	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &CheckCommand{
						client: http.DefaultClient,
					}
				},
			}
		})
		return pool
	}
}

// CheckCommand is a reusable command for making HTTP GET requests.
type CheckCommand struct {
	client *http.Client // HTTP client used for performing requests.
	url    string       // url is the target endpoint for the HTTP GET request.
}

// GetCheckCommand retrieves CheckCommand from the pool.
func GetCheckCommand() *CheckCommand {
	return checkCommandPool().Get().(*CheckCommand)
}

// SetClient configures the CheckCommand with a custom *http.Client and URL.
func (c *CheckCommand) SetClient(client *http.Client, url string) *CheckCommand {
	c.client = client
	c.url = url
	return c
}

// Reset clears the client and URL preparing the CheckCommand to be reused later.
func (c *CheckCommand) Reset() *CheckCommand {
	c.client = nil
	c.url = ""
	return c
}

// Release puts the current CheckCommand back into the pool.
func (c *CheckCommand) Release() {
	checkCommandPool().Put(c.Reset())
}

// Execute performs HTTP GET request to the configured URL using the command's *http.Client.
func (c *CheckCommand) Execute(ctx context.Context) (result string, err error) {
	var (
		request  *http.Request
		response *http.Response
		body     []byte
	)

	if request, err = http.NewRequestWithContext(ctx, http.MethodGet, c.url, http.NoBody); err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	// Perform the request
	if response, err = c.client.Do(request); err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			fmt.Printf("close response body: %v", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status: %s", response.Status)
	}
	if body, err = io.ReadAll(response.Body); err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}
	return string(body), nil
}
