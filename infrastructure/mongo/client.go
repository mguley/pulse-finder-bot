package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient creates a new MongoDB client and establishes a connection.
func NewMongoClient(uri string) (*mongo.Client, error) {
	// Configure the client options using the provided URI
	clientOptions := options.Client().ApplyURI(uri)

	// Establish a connection to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to mongo, %w", err)
	}
	return client, nil
}
