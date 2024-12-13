package url

import (
	"context"
	"domain/url/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestRepository_Save validates the Save method of the URL repository.
func TestRepository_Save(t *testing.T) {
	container := SetupTestContainer(t)
	repo := container.UrlRepository.Get()

	ctx := context.Background()
	url := &entity.Url{
		ID:      primitive.NewObjectID(),
		Address: "https://example.com",
		Status:  "pending",
	}

	// Save the URL entity
	err := repo.Save(ctx, url)
	require.NoError(t, err, "Failed to save URL entity")

	// Verify the entity exists in the database
	var result []*entity.Url
	result, err = repo.FetchBatch(ctx, "pending", 1)
	require.NoError(t, err, "Failed to fetch batch")
	require.Len(t, result, 1, "Unexpected number of results")
	assert.Equal(t, url.Address, result[0].Address, "Address is not as expected")
	assert.Equal(t, url.Status, result[0].Status, "Status is not as expected")
	assert.Equal(t, url.ID, result[0].ID, "ID is not as expected")
}

// TestRepository_FetchBatch validates the FetchBatch method of the URL repository.
func TestRepository_FetchBatch(t *testing.T) {
	container := SetupTestContainer(t)
	repo := container.UrlRepository.Get()

	ctx := context.Background()
	testData := []*entity.Url{
		{ID: primitive.NewObjectID(), Address: "https://example1.com", Status: "pending"},
		{ID: primitive.NewObjectID(), Address: "https://example2.com", Status: "pending"},
		{ID: primitive.NewObjectID(), Address: "https://example3.com", Status: "completed"},
	}

	// Seed the database with test data
	for _, url := range testData {
		err := repo.Save(ctx, url)
		require.NoError(t, err, "Failed to save URL entity")
	}

	// Fetch a batch of URLs with status "pending"
	results, err := repo.FetchBatch(ctx, "pending", 2)
	require.NoError(t, err, "Failed to fetch batch")
	assert.Len(t, results, 2, "Unexpected number of results")
	for _, result := range results {
		assert.Equal(t, "pending", result.Status, "Status is not as expected")
	}
}

// TestRepository_UpdateStatus validates the UpdateStatus method of the URL repository.
func TestRepository_UpdateStatus(t *testing.T) {
	container := SetupTestContainer(t)
	repo := container.UrlRepository.Get()

	ctx := context.Background()
	// Seed the database with a test entity
	testUrl := &entity.Url{
		ID:      primitive.NewObjectID(),
		Address: "https://example.com",
		Status:  "pending",
	}
	err := repo.Save(ctx, testUrl)
	require.NoError(t, err, "Failed to save URL entity")

	// Update the status
	newStatus := "completed"
	processedTime := time.Now()
	err = repo.UpdateStatus(ctx, testUrl.ID.Hex(), newStatus, &processedTime)
	require.NoError(t, err, "Failed to update status")

	// Verify the update
	results, err := repo.FetchBatch(ctx, "completed", 1)
	require.NoError(t, err, "Failed to fetch batch")
	assert.Len(t, results, 1, "Unexpected number of results")
	assert.Equal(t, newStatus, results[0].Status, "Status is not as expected")
	assert.Equal(t, testUrl.ID, results[0].ID, "ID is not as expected")
}
