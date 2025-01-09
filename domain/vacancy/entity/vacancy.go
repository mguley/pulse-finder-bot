package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Vacancy represents a job vacancy entity.
// It captures details such as the title, company, description, posting date, and location.
type Vacancy struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`        // Unique identifier for the URL
	Title       string             `bson:"title" json:"title"`             // The title of the job vacancy.
	Company     string             `bson:"company" json:"company"`         // The company offering the job.
	Description string             `bson:"description" json:"description"` // A detailed description of the job vacancy.
	PostedAt    time.Time          `bson:"posted_at" json:"postedAt"`      // The date and time when it was posted.
	Location    string             `bson:"location" json:"location"`       // The location of the job vacancy.
	SentAt      time.Time          `bson:"sent_at" json:"sentAt"`          // The timestamp when the vacancy was sent.
}
