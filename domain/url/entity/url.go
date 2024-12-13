package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Url represents a URL entity.
// It captures information about the URL to be processed, its status, and metadata.
type Url struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`    // Unique identifier for the URL (MongoDB ObjectID).
	Address   string             `bson:"address" json:"address"`     // The URL address to be processed.
	Status    string             `bson:"status" json:"status"`       // Current processing status of the URL.
	Processed time.Time          `bson:"processed" json:"processed"` // Timestamp of when the URL was processed.
}
