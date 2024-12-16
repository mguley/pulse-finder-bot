package sitemap

import (
	"application/config"
	"application/dependency"
	"application/proxy/services"
	"domain/url/repository"
	"domain/useragent"
	"fmt"
	httpClient "infrastructure/http/client"
	infraMongo "infrastructure/mongo"
	"infrastructure/proxy/client"
	"infrastructure/proxy/client/agent"
	"infrastructure/url"
	"infrastructure/url/sitemap/fetcher"
	"infrastructure/url/sitemap/notifier"
	"infrastructure/url/sitemap/parser"
	sitemapRepository "infrastructure/url/sitemap/repository"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config            dependency.LazyDependency[*config.Config]
	UserAgent         dependency.LazyDependency[useragent.Generator]
	Socks5Client      dependency.LazyDependency[*client.Socks5Client]
	HttpFactory       dependency.LazyDependency[*httpClient.Factory]
	ProxyService      dependency.LazyDependency[*services.Service]
	SitemapFetcher    dependency.LazyDependency[*fetcher.Service]
	SitemapNotifier   dependency.LazyDependency[*notifier.Service]
	SitemapParser     dependency.LazyDependency[*parser.Service]
	MongoClient       dependency.LazyDependency[*mongo.Client]
	UrlRepository     dependency.LazyDependency[repository.UrlRepository]
	SitemapRepository dependency.LazyDependency[*sitemapRepository.Service]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.Config = dependency.LazyDependency[*config.Config]{
		InitFunc: config.LoadConfig,
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
	c.ProxyService = dependency.LazyDependency[*services.Service]{
		InitFunc: func() *services.Service {
			factory := c.HttpFactory.Get()
			host := c.Config.Get().Proxy.Host
			port := c.Config.Get().Proxy.Port
			timeout := 10 * time.Second
			return services.NewService(factory, host, port, timeout)
		},
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
			cfg := c.Config.Get()
			mongoClient := c.MongoClient.Get()
			collection := mongoClient.Database(cfg.Mongo.DB).Collection(cfg.Mongo.UrlsCollection)
			return url.NewRepository(mongoClient, collection)
		},
	}

	// Sitemap services
	c.SitemapFetcher = dependency.LazyDependency[*fetcher.Service]{
		InitFunc: func() *fetcher.Service {
			proxyClient := func() (*http.Client, error) {
				return c.ProxyService.Get().HttpClient()
			}
			return fetcher.NewService(proxyClient)
		},
	}
	c.SitemapNotifier = dependency.LazyDependency[*notifier.Service]{
		InitFunc: func() *notifier.Service {
			return notifier.NewService(c.Config.Get().Proxy.PingUrl)
		},
	}
	c.SitemapParser = dependency.LazyDependency[*parser.Service]{
		InitFunc: parser.NewService,
	}
	c.SitemapRepository = dependency.LazyDependency[*sitemapRepository.Service]{
		InitFunc: func() *sitemapRepository.Service {
			return sitemapRepository.NewService(c.UrlRepository.Get())
		},
	}

	return c
}
