package repository

import (
	"context"
	"domain/vacancy/entity"
)

// VacancyRepository defines the interface for interacting with vacancy entities in the persistence layer.
type VacancyRepository interface {
	// Save persists a new vacancy into the data source.
	// Returns an error if the operation fails.
	Save(ctx context.Context, vacancy *entity.Vacancy) error

	// Update updates an existing vacancy by ID.
	// Returns an error if the operation fails.
	Update(ctx context.Context, vacancy *entity.Vacancy) error

	// Fetch retrieves a list of vacancies with optional filters.
	// Returns an error if the operation fails.
	Fetch(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*entity.Vacancy, error)

	// FetchBatch retrieves a batch of vacancies where the SentAt field is not set.
	// Returns a slice of Vacancy entities matching the criteria.
	FetchBatch(ctx context.Context, limit int) ([]*entity.Vacancy, error)

	// FindByID retrieves a vacancy by its ID.
	// Returns an error if the operation fails.
	FindByID(ctx context.Context, id string) (*entity.Vacancy, error)
}
