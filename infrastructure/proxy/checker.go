package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var checkerPool = newCheckerPool()

// newCheckerPool creates and initializes a singleton sync.Pool for managing reusable Checker objects.
// This uses a once.Do mechanism to ensure the pool is only created once during runtime.
func newCheckerPool() func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool

	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &Checker{}
				},
			}
		})
		return pool
	}
}

// Checker is a reusable struct that performs HTTP requests.
type Checker struct {
	client *http.Client // The HTTP client used for performing requests.
}

// SetClient sets a custom HTTP client for the Checker instance.
func (c *Checker) SetClient(client *http.Client) *Checker {
	c.client = client
	return c
}

// GetInfo fetches content from the provided URL using the configured HTTP client.
func (c *Checker) GetInfo(ctx context.Context, url string) (string, error) {
	// If no HTTP client is set, use the default HTTP client.
	if c.client == nil {
		c.client = http.DefaultClient
	}

	// Create a new HTTP request.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	// Perform the HTTP request.
	res, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			fmt.Printf("close response body error: %v", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status: %d", res.StatusCode)
	}

	// Read and return the response body.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}
	return string(body), nil
}

// Reset clears the state of the Checker object, preparing it for reuse.
func (c *Checker) Reset() *Checker {
	c.client = nil
	return c
}

// Release returns the Checker instance to the pool after resetting it.
func (c *Checker) Release() {
	checkerPool().Put(c.Reset())
}

// GetChecker retrieves a Checker instance from the pool, ensuring it is reset before use.
func GetChecker() *Checker {
	return checkerPool().Get().(*Checker).Reset()
}
