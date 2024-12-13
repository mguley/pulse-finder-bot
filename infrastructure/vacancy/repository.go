package vacancy

import (
	"context"
	"domain/vacancy/entity"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository provides a MongoDB based implementation for managing vacancy entities.
type Repository struct {
	client     *mongo.Client     // MongoDB client instance.
	collection *mongo.Collection // MongoDB collection for storing URL entities.
}

// NewRepository creates a new Repository instance.
func NewRepository(client *mongo.Client, collection *mongo.Collection) *Repository {
	return &Repository{client: client, collection: collection}
}

// Save persists a new vacancy entity into the MongoDB collection.
func (r *Repository) Save(ctx context.Context, vacancy *entity.Vacancy) error {
	if vacancy.ID.IsZero() {
		vacancy.ID = primitive.NewObjectID()
	}
	if _, err := r.collection.InsertOne(ctx, vacancy); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// Fetch retrieves a list of vacancies with optional filters and pagination.
func (r *Repository) Fetch(
	ctx context.Context,
	filters map[string]interface{},
	limit, offset int,
) ([]*entity.Vacancy, error) {
	filter := bson.M(filters)
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "posted_at", Value: -1}}) // Default sort by most recent

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			fmt.Println("cursor.Close", err)
		}
	}()

	var list []*entity.Vacancy
	if err = cursor.All(ctx, &list); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return list, nil
}

// FindByID retrieves a vacancy entity by its ID from the MongoDB collection.
func (r *Repository) FindByID(ctx context.Context, id string) (*entity.Vacancy, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	filter := bson.M{"_id": oid}
	vacancy := &entity.Vacancy{}
	if err = r.collection.FindOne(ctx, filter).Decode(vacancy); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return vacancy, nil
}
