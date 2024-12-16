package repository

import (
	"context"
	"domain/url/entity"
	"domain/url/repository"
	"fmt"
	"time"
)

// Service handles operations for saving URLs to the data source.
type Service struct {
	urlRepository repository.UrlRepository // The repository instance for managing URL entities.
}

// NewService creates and returns a new Repository service instance.
func NewService(r repository.UrlRepository) *Service {
	return &Service{urlRepository: r}
}

// SaveUrls saves a batch of URLs to the data source.
func (s *Service) SaveUrls(ctx context.Context, urls []string) error {
	for _, url := range urls {
		item := entity.Url{
			Address:   url,
			Status:    "pending",
			Processed: time.Time{}, // Not yet processed.
		}

		if err := s.urlRepository.Save(ctx, &item); err != nil {
			return fmt.Errorf("save URL %s: %w", url, err)
		}
	}

	return nil
}
