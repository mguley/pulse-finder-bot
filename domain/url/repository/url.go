package repository

import (
	"context"
	"domain/url/entity"
	"time"
)

// UrlRepository defines the interface for interacting with URL entities in the persistence layer.
type UrlRepository interface {
	// Save persists a new URL entity into the data source.
	// Returns an error if the operation fails.
	Save(ctx context.Context, url *entity.Url) error

	// FetchBatch retrieves a batch of URLs with the specified status.
	// Returns a slice of URL entities matching the criteria.
	FetchBatch(ctx context.Context, status string, limit int) ([]*entity.Url, error)

	// UpdateStatus updates the status of URL entity in the data source.
	// Returns an error if the operation fails.
	UpdateStatus(ctx context.Context, id, status string, processedTime *time.Time) error
}
