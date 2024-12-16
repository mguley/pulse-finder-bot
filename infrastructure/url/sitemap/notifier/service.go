package notifier

import (
	"context"
	"fmt"
	"infrastructure/proxy"
	"net/http"
	"time"
)

// Service handles notification tasks, such as logging proxy IP information.
type Service struct {
	url string // The URL to use for fetching IP information.
}

// NewService creates and returns a new Notifier service instance.
func NewService(url string) *Service {
	return &Service{url: url}
}

// Notify retrieves and logs the proxy's IP address using the provided HTTP client.
func (s *Service) Notify(client *http.Client) error {
	checker := proxy.GetChecker().SetClient(client)
	defer checker.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ip, err := checker.GetInfo(ctx, s.url)
	if err != nil {
		return fmt.Errorf("get info: %w", err)
	}

	fmt.Printf("We use: %s\n", ip)
	return nil
}
