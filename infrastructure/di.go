package infrastructure

import (
	"application/config"
	"application/dependency"
	"domain/useragent"
	"fmt"
	httpClient "infrastructure/http/client"
	"infrastructure/proxy/client"
	"infrastructure/proxy/client/agent"
	"infrastructure/proxy/port"
	"time"
)

// Container provides a lazily initialized set of dependencies for the infrastructure layer.
type Container struct {
	ProxyConnection dependency.LazyDependency[*port.Connection]
	UserAgent       dependency.LazyDependency[useragent.Generator]
	Socks5Client    dependency.LazyDependency[*client.Socks5Client]
	HttpFactory     dependency.LazyDependency[*httpClient.Factory]
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

	return c
}
