package vacancy

import (
	"application/config"
	"application/dependency"
	vacancyRepo "domain/vacancy/repository"
	"fmt"
	infraMongo "infrastructure/mongo"
	"infrastructure/vacancy"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config            dependency.LazyDependency[*config.Config]
	MongoClient       dependency.LazyDependency[*mongo.Client]
	VacancyRepository dependency.LazyDependency[vacancyRepo.VacancyRepository]
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
	c.VacancyRepository = dependency.LazyDependency[vacancyRepo.VacancyRepository]{
		InitFunc: func() vacancyRepo.VacancyRepository {
			mongoClient := c.MongoClient.Get()
			cfg := c.Config.Get()
			collection := mongoClient.Database(cfg.Mongo.DB).Collection(cfg.Mongo.VacancyCollection)
			return vacancy.NewRepository(mongoClient, collection)
		},
	}

	return c
}
