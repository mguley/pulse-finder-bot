package infrastructure

import (
	"application/config"
	"application/dependency"
	"domain/url/repository"
	"domain/useragent"
	"fmt"
	httpClient "infrastructure/http/client"
	infraMongo "infrastructure/mongo"
	"infrastructure/proxy/client"
	"infrastructure/proxy/client/agent"
	"infrastructure/proxy/port"
	"infrastructure/url"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Container provides a lazily initialized set of dependencies for the infrastructure layer.
type Container struct {
	ProxyConnection dependency.LazyDependency[*port.Connection]
	UserAgent       dependency.LazyDependency[useragent.Generator]
	Socks5Client    dependency.LazyDependency[*client.Socks5Client]
	HttpFactory     dependency.LazyDependency[*httpClient.Factory]
	MongoClient     dependency.LazyDependency[*mongo.Client]
	UrlRepository   dependency.LazyDependency[repository.UrlRepository]
}

// NewContainer initializes and returns a new Container with lazy dependencies for the infrastructure layer.
func NewContainer(cfg *config.Config) *Container {
	c := &Container{}

	c.ProxyConnection = dependency.LazyDependency[*port.Connection]{
		InitFunc: func() *port.Connection {
			address := fmt.Sprintf("%s:%s", cfg.Proxy.Host, cfg.Proxy.ControlPort)
			timeout := 10 * time.Second
			return port.NewConnection(address, cfg.Proxy.ControlPassword, timeout)
		},
	}
	c.UserAgent = dependency.LazyDependency[useragent.Generator]{
		InitFunc: func() useragent.Generator {
			return agent.NewChromeUserAgentGenerator()
		},
	}
	c.Socks5Client = dependency.LazyDependency[*client.Socks5Client]{
		InitFunc: func() *client.Socks5Client {
			return client.NewSocks5Client(c.UserAgent.Get())
		},
	}
	c.HttpFactory = dependency.LazyDependency[*httpClient.Factory]{
		InitFunc: func() *httpClient.Factory {
			return httpClient.NewFactory(c.UserAgent.Get(), c.Socks5Client.Get())
		},
	}
	c.MongoClient = dependency.LazyDependency[*mongo.Client]{
		InitFunc: func() *mongo.Client {
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
			collection := mongoClient.Database(cfg.Mongo.DB).Collection(cfg.Mongo.UrlsCollection)
			return url.NewRepository(mongoClient, collection)
		},
	}

	return c
}
