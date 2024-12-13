package url

import (
	"application/config"
	"application/dependency"
	"domain/url/repository"
	"fmt"
	infraMongo "infrastructure/mongo"
	"infrastructure/url"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config        dependency.LazyDependency[*config.Config]
	MongoClient   dependency.LazyDependency[*mongo.Client]
	UrlRepository dependency.LazyDependency[repository.UrlRepository]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.Config = dependency.LazyDependency[*config.Config]{
		InitFunc: config.LoadConfig,
	}
	c.MongoClient = dependency.LazyDependency[*mongo.Client]{
		InitFunc: func() *mongo.Client {
			cfg := c.Config.Get()
			uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.Mongo.User, cfg.Mongo.Pass, cfg.Mongo.Host, cfg.Mongo.Port)
			mongoClient, err := infraMongo.NewMongoClient(uri)
			if err != nil {
				log.Fatalf("mongo client error: %v", err)
			}
			return mongoClient
		},
	}
	c.UrlRepository = dependency.LazyDependency[repository.UrlRepository]{
		InitFunc: func() repository.UrlRepository {
			mongoClient := c.MongoClient.Get()
			cfg := c.Config.Get()
			collection := mongoClient.Database(cfg.Mongo.DB).Collection(cfg.Mongo.UrlsCollection)
			return url.NewRepository(mongoClient, collection)
		},
	}

	return c
}
