package beta

import (
	"application/config"
	"application/dependency"
	"application/proxy/services"
	"domain/html"
	"domain/useragent"
	htmlBeta "infrastructure/html/source/beta"
	httpClient "infrastructure/http/client"
	"infrastructure/proxy/client"
	"infrastructure/proxy/client/agent"
	"log"
	"time"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config          dependency.LazyDependency[*config.Config]
	UserAgent       dependency.LazyDependency[useragent.Generator]
	Socks5Client    dependency.LazyDependency[*client.Socks5Client]
	HttpFactory     dependency.LazyDependency[*httpClient.Factory]
	ProxyService    dependency.LazyDependency[*services.Service]
	BetaHtmlFetcher dependency.LazyDependency[html.Fetcher]
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
	c.BetaHtmlFetcher = dependency.LazyDependency[html.Fetcher]{
		InitFunc: func() html.Fetcher {
			maxBodySize := int64(10 * 1024 * 1024)
			fetcher, err := htmlBeta.NewFetcher(c.ProxyService.Get(), maxBodySize)
			if err != nil {
				log.Fatalf("failed to init fetcher: %v", err)
			}
			return fetcher
		},
	}

	return c
}
