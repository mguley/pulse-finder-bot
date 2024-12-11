package services

import (
	"application/config"
	"application/dependency"
	"application/proxy/services"
	"domain/useragent"
	httpClient "infrastructure/http/client"
	"infrastructure/proxy/client"
	"infrastructure/proxy/client/agent"
	"time"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config       dependency.LazyDependency[*config.Config]
	UserAgent    dependency.LazyDependency[useragent.Generator]
	Socks5Client dependency.LazyDependency[*client.Socks5Client]
	HttpFactory  dependency.LazyDependency[*httpClient.Factory]
	ProxyService dependency.LazyDependency[*services.Service]
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

	// Proxy services
	c.ProxyService = dependency.LazyDependency[*services.Service]{
		InitFunc: func() *services.Service {
			factory := c.HttpFactory.Get()
			host := c.Config.Get().Proxy.Host
			port := c.Config.Get().Proxy.Port
			timeout := 10 * time.Second
			return services.NewService(factory, host, port, timeout)
		},
	}

	return c
}
