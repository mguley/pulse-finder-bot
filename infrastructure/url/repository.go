package url

import (
	"context"
	"domain/url/entity"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository provides a MongoDB based implementation for managing URL entities.
type Repository struct {
	client     *mongo.Client     // MongoDB client instance.
	collection *mongo.Collection // MongoDB collection for storing URL entities.
}

// NewRepository creates a new Repository instance.
func NewRepository(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{client: client, collection: collection}
}

// Save persists a new URL entity into the MongoDB collection.
func (r *Repository) Save(ctx context.Context, url *entity.Url) error {
	if url.ID.IsZero() {
		url.ID = primitive.NewObjectID()
	}
	if _, err := r.collection.InsertOne(ctx, url); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// FetchBatch retrieves a batch of URLs with the specified status from MongoDB.
func (r *Repository) FetchBatch(ctx context.Context, status string, limit int) ([]*entity.Url, error) {
	filter := bson.M{"status": status}
	opt := options.Find().SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opt)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			fmt.Println("cursor.Close", err)
		}
	}()

	var urls []*entity.Url
	if err = cursor.All(ctx, &urls); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return urls, nil
}

// UpdateStatus updates the status of URL entity in the MongoDB collection.
func (r *Repository) UpdateStatus(ctx context.Context, id, status string, processedTime *time.Time) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	update := bson.M{"$set": bson.M{
		"status": status,
	}}
	if processedTime != nil {
		update["$set"].(bson.M)["processed"] = *processedTime
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return fmt.Errorf("update document: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no document found with the id %s", id)
	}
	return nil
}
