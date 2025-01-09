package cron

import (
	"context"
	"domain/vacancy/entity"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestCronScheduler_Integration ensures that the CronScheduler processes correctly unsent items (SentAt=zero).
func TestCronScheduler_Integration(t *testing.T) {
	container := SetupTestEnvironment(t)

	scheduler := container.CronScheduler.Get()
	repo := container.VacancyRepository.Get()
	ctx := context.Background()

	// Seed the database with test data
	now := time.Now()
	testData := []*entity.Vacancy{
		{ID: primitive.NewObjectID(), Title: "Job 1", Company: "Company A", Description: "Desc 1", PostedAt: now, Location: "Location 1", SentAt: time.Time{}},
		{ID: primitive.NewObjectID(), Title: "Job 2", Company: "Company B", Description: "Desc 2", PostedAt: now.Add(-time.Hour), Location: "Location 2", SentAt: time.Time{}},
		{ID: primitive.NewObjectID(), Title: "Job 3", Company: "Company C", Description: "Desc 3", PostedAt: now.Add(-2 * time.Hour), Location: "Location 3", SentAt: now.Add(-1 * time.Hour)},
		{ID: primitive.NewObjectID(), Title: "Job 4", Company: "Company D", Description: "Desc 4", PostedAt: now.Add(-3 * time.Hour), Location: "Location 4", SentAt: time.Time{}},
		{ID: primitive.NewObjectID(), Title: "Job 5", Company: "Company E", Description: "Desc 5", PostedAt: now.Add(-4 * time.Hour), Location: "Location 5", SentAt: now.Add(-2 * time.Hour)},
	}

	for _, vacancy := range testData {
		err := repo.Save(ctx, vacancy)
		require.NoError(t, err, "Failed to save vacancy")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the CronScheduler
	scheduler.Start(ctx)

	// Give the scheduler a little time to process, we set tickerTime = 5 seconds (cron/di.go)
	fmt.Println("sleeping...")
	time.Sleep(10 * time.Second)

	// Verify that all inserted items have a non-zero SentAt
	items, err := repo.Fetch(ctx, bson.M{}, 5, 0)
	require.NoError(t, err, "Failed to fetch items after cron run")
	require.NotEmpty(t, items, "We should have found at least one item")

	for _, v := range items {
		assert.False(t, v.SentAt.IsZero(),
			"Vacancy %s should have a non-zero SentAt after processing", v.ID.Hex())
	}
	defer scheduler.Stop()
}
