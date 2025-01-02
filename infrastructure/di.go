package infrastructure

import (
	"application/config"
	"application/dependency"
	"domain/url/repository"
	"domain/useragent"
	vacancyRepo "domain/vacancy/repository"
	"fmt"
	authClient "infrastructure/grpc/auth/client"
	httpClient "infrastructure/http/client"
	infraMongo "infrastructure/mongo"
	proxyClient "infrastructure/proxy/client"
	proxyAgent "infrastructure/proxy/client/agent"
	proxyPort "infrastructure/proxy/port"
	"infrastructure/url"
	"infrastructure/vacancy"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Container provides a lazily initialized set of dependencies for the infrastructure layer.
type Container struct {
	ProxyConnection   dependency.LazyDependency[*proxyPort.Connection]
	UserAgent         dependency.LazyDependency[useragent.Generator]
	Socks5Client      dependency.LazyDependency[*proxyClient.Socks5Client]
	HttpFactory       dependency.LazyDependency[*httpClient.Factory]
	MongoClient       dependency.LazyDependency[*mongo.Client]
	UrlRepository     dependency.LazyDependency[repository.UrlRepository]
	VacancyRepository dependency.LazyDependency[vacancyRepo.VacancyRepository]
	AuthClient        dependency.LazyDependency[*authClient.AuthClient]
}

// NewContainer initializes and returns a new Container with lazy dependencies for the infrastructure layer.
func NewContainer(cfg *config.Config) *Container {
	c := &Container{}

	c.ProxyConnection = dependency.LazyDependency[*proxyPort.Connection]{
		InitFunc: func() *proxyPort.Connection {
			address := fmt.Sprintf("%s:%s", cfg.Proxy.Host, cfg.Proxy.ControlPort)
			timeout := 10 * time.Second
			return proxyPort.NewConnection(address, cfg.Proxy.ControlPassword, timeout)
		},
	}
	c.UserAgent = dependency.LazyDependency[useragent.Generator]{
		InitFunc: func() useragent.Generator {
			return proxyAgent.NewChromeUserAgentGenerator()
		},
	}
	c.Socks5Client = dependency.LazyDependency[*proxyClient.Socks5Client]{
		InitFunc: func() *proxyClient.Socks5Client {
			return proxyClient.NewSocks5Client(c.UserAgent.Get())
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
			uri = uri + "/" + cfg.Mongo.DB + "?authSource=admin"
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
	c.VacancyRepository = dependency.LazyDependency[vacancyRepo.VacancyRepository]{
		InitFunc: func() vacancyRepo.VacancyRepository {
			mongoClient := c.MongoClient.Get()
			collection := mongoClient.Database(cfg.Mongo.DB).Collection(cfg.Mongo.VacancyCollection)
			return vacancy.NewRepository(mongoClient, collection)
		},
	}
	c.AuthClient = dependency.LazyDependency[*authClient.AuthClient]{
		InitFunc: func() *authClient.AuthClient {
			env := cfg.Env
			address := cfg.AuthServer.Address
			client, err := authClient.NewAuthClient(env, address)
			if err != nil {
				log.Fatalf("auth client error: %v", err)
			}
			return client
		},
	}

	return c
}
