package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var checkerPool = newCheckerPool()

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

// Checker is a reusable struct for making HTTP requests.
type Checker struct {
	client *http.Client
}

// SetClient sets HTTP client for the Checker.
func (c *Checker) SetClient(client *http.Client) *Checker {
	c.client = client
	return c
}

// GetInfo fetches content from the given URL using the configured HTTP client.
func (c *Checker) GetInfo(ctx context.Context, url string) (string, error) {
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

// Reset clears a Checker state for reuse.
func (c *Checker) Reset() *Checker {
	c.client = nil
	return c
}

// Release puts a Checker instance back into the pool.
func (c *Checker) Release() {
	checkerPool().Put(c.Reset())
}

// GetChecker retrieves a Checker instance from the pool.
func GetChecker() *Checker {
	return checkerPool().Get().(*Checker).Reset()
}
