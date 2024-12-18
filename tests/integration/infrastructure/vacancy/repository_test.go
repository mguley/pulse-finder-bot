package vacancy

import (
	"context"
	"domain/vacancy/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestRepository_Save validates the Save method of the VacancyRepository.
func TestRepository_Save(t *testing.T) {
	container := SetupTestContainer(t)
	repo := container.VacancyRepository.Get()

	ctx := context.Background()
	vacancy := &entity.Vacancy{
		ID:          primitive.NewObjectID(),
		Title:       "Software Engineer",
		Company:     "TechCorp",
		Description: "Develop and maintain software applications.",
		PostedAt:    time.Now(),
		Location:    "New York",
	}

	// Save the vacancy entity
	err := repo.Save(ctx, vacancy)
	require.NoError(t, err, "Failed to save vacancy entity")

	// Verify the entity exists in the database
	result, err := repo.FindByID(ctx, vacancy.ID.Hex())
	require.NoError(t, err, "Failed to find vacancy by ID")
	assert.Equal(t, vacancy.Title, result.Title, "Title is not as expected")
	assert.Equal(t, vacancy.Company, result.Company, "Company is not as expected")
	assert.Equal(t, vacancy.Description, result.Description, "Description is not as expected")
	assert.Equal(t, vacancy.Location, result.Location, "Location is not as expected")
	assert.WithinDuration(t, vacancy.PostedAt, result.PostedAt, time.Second, "PostedAt timestamp mismatch")
}

// TestRepository_Fetch validates the Fetch method of the VacancyRepository.
func TestRepository_Fetch(t *testing.T) {
	container := SetupTestContainer(t)
	repo := container.VacancyRepository.Get()

	ctx := context.Background()
	// Seed the database with test data
	testData := []*entity.Vacancy{
		{ID: primitive.NewObjectID(), Title: "Backend Developer", Company: "TechCorp", Description: "Work on backend systems.", PostedAt: time.Now().Add(-48 * time.Hour), Location: "San Francisco"},
		{ID: primitive.NewObjectID(), Title: "Frontend Developer", Company: "WebDesigns", Description: "Work on frontend systems.", PostedAt: time.Now().Add(-24 * time.Hour), Location: "New York"},
		{ID: primitive.NewObjectID(), Title: "Data Scientist", Company: "DataWorks", Description: "Analyze data.", PostedAt: time.Now().Add(-72 * time.Hour), Location: "Boston"},
	}

	for _, vacancy := range testData {
		err := repo.Save(ctx, vacancy)
		require.NoError(t, err, "Failed to save vacancy")
	}

	// Fetch vacancies with specific filters
	filters := map[string]interface{}{"company": "TechCorp"}
	results, err := repo.Fetch(ctx, filters, 2, 0)
	require.NoError(t, err, "Failed to fetch vacancies")
	assert.Len(t, results, 1, "Unexpected number of results")
	assert.Equal(t, "Backend Developer", results[0].Title, "Title is not as expected")
	assert.Equal(t, "TechCorp", results[0].Company, "Company is not as expected")
}

// TestRepository_FindByID validates the FindByID method of the VacancyRepository.
func TestRepository_FindByID(t *testing.T) {
	container := SetupTestContainer(t)
	repo := container.VacancyRepository.Get()

	ctx := context.Background()
	// Seed the database with a test entity
	testVacancy := &entity.Vacancy{
		ID:          primitive.NewObjectID(),
		Title:       "Cloud Engineer",
		Company:     "CloudTech",
		Description: "Manage cloud infrastructure.",
		PostedAt:    time.Now(),
		Location:    "Remote",
	}
	err := repo.Save(ctx, testVacancy)
	require.NoError(t, err, "Failed to save vacancy entity")

	// Find the vacancy by ID
	result, err := repo.FindByID(ctx, testVacancy.ID.Hex())
	require.NoError(t, err, "Failed to find vacancy by ID")
	assert.Equal(t, testVacancy.Title, result.Title, "Title is not as expected")
	assert.Equal(t, testVacancy.Company, result.Company, "Company is not as expected")
	assert.Equal(t, testVacancy.Location, result.Location, "Location is not as expected")
	assert.WithinDuration(t, testVacancy.PostedAt, result.PostedAt, time.Second, "PostedAt timestamp mismatch")
}