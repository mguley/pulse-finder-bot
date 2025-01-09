package cron

import (
	"application/config"
	"application/dependency"
	appScheduler "application/scheduler"
	"domain/scheduler"
	vacancyRepo "domain/vacancy/repository"
	"fmt"
	infraMongo "infrastructure/mongo"
	"infrastructure/vacancy"
	"log"
	authClient "tests/integration/application/scheduler/cron/auth/client"
	vacancyClient "tests/integration/application/scheduler/cron/vacancy/client"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config                 dependency.LazyDependency[*config.Config]
	MongoClient            dependency.LazyDependency[*mongo.Client]
	VacancyRepository      dependency.LazyDependency[vacancyRepo.VacancyRepository]
	AuthClientContainer    dependency.LazyDependency[*authClient.TestContainer]
	VacancyClientContainer dependency.LazyDependency[*vacancyClient.TestContainer]
	CronScheduler          dependency.LazyDependency[scheduler.Scheduler]
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
			uri = uri + "/" + cfg.Mongo.DB + "?authSource=admin"
			mongoClient, err := infraMongo.NewMongoClient(uri)
			if err != nil {
				log.Fatalf("mongo client error: %v", err)
			}
			return mongoClient
		},
	}
	c.VacancyRepository = dependency.LazyDependency[vacancyRepo.VacancyRepository]{
		InitFunc: func() vacancyRepo.VacancyRepository {
			cfg := c.Config.Get()
			mongoClient := c.MongoClient.Get()
			collection := mongoClient.Database(cfg.Mongo.DB).Collection(cfg.Mongo.VacancyCollection)
			return vacancy.NewRepository(mongoClient, collection)
		},
	}
	c.AuthClientContainer = dependency.LazyDependency[*authClient.TestContainer]{
		InitFunc: authClient.NewTestContainer,
	}
	c.VacancyClientContainer = dependency.LazyDependency[*vacancyClient.TestContainer]{
		InitFunc: vacancyClient.NewTestContainer,
	}
	c.CronScheduler = dependency.LazyDependency[scheduler.Scheduler]{
		InitFunc: func() scheduler.Scheduler {
			cfg := c.Config.Get()
			repo := c.VacancyRepository.Get()
			aClient := c.AuthClientContainer.Get().AuthClient.Get()
			vClient := c.VacancyClientContainer.Get().VacancyClient.Get()
			batchSize := 5
			issuer := cfg.AuthServer.Issuer
			scope := []string{"write"}
			tickerTime := 5 * time.Second
			return appScheduler.NewCronScheduler(repo, aClient, vClient, batchSize, issuer, scope, tickerTime)
		},
	}

	return c
}
