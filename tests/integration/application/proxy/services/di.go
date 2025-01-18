package services

import (
	"application/config"
	"application/dependency"
	"application/proxy/circuit"
	"application/proxy/commands"
	"application/proxy/commands/control"
	"application/proxy/services"
	"application/proxy/strategies"
	"domain/useragent"
	"fmt"
	httpClient "infrastructure/http/client"
	"infrastructure/proxy/client"
	"infrastructure/proxy/client/agent"
	"infrastructure/proxy/port"
	"time"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config              dependency.LazyDependency[*config.Config]
	UserAgent           dependency.LazyDependency[useragent.Generator]
	Socks5Client        dependency.LazyDependency[*client.Socks5Client]
	HttpFactory         dependency.LazyDependency[*httpClient.Factory]
	ProxyService        dependency.LazyDependency[*services.Service]
	ProxyConnection     dependency.LazyDependency[*port.Connection]
	AuthenticateCommand dependency.LazyDependency[*control.AuthenticateCommand]
	SignalCommand       dependency.LazyDependency[*control.SignalCommand]
	StatusCommand       dependency.LazyDependency[*commands.StatusCommand]
	RetryStrategy       dependency.LazyDependency[strategies.RetryStrategy]
	IdentityService     dependency.LazyDependency[*services.Identity]
	CircuitManager      dependency.LazyDependency[*circuit.Manager]
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
	c.ProxyConnection = dependency.LazyDependency[*port.Connection]{
		InitFunc: func() *port.Connection {
			address := fmt.Sprintf("%s:%s", c.Config.Get().Proxy.Host, c.Config.Get().Proxy.ControlPort)
			timeout := time.Duration(10) * time.Second
			return port.NewConnection(address, c.Config.Get().Proxy.ControlPassword, timeout)
		},
	}
	c.RetryStrategy = dependency.LazyDependency[strategies.RetryStrategy]{
		InitFunc: func() strategies.RetryStrategy {
			baseDelay := time.Duration(5) * time.Second
			maxDelay := time.Duration(30) * time.Second
			maxAttempts := 5
			multiplier := 2.0
			return strategies.NewExponentialBackoffStrategy(baseDelay, maxDelay, maxAttempts, multiplier)
		},
	}

	// Proxy commands
	c.AuthenticateCommand = dependency.LazyDependency[*control.AuthenticateCommand]{
		InitFunc: func() *control.AuthenticateCommand {
			return control.NewAuthenticateCommand(c.ProxyConnection.Get())
		},
	}
	c.SignalCommand = dependency.LazyDependency[*control.SignalCommand]{
		InitFunc: func() *control.SignalCommand {
			return control.NewSignalCommand(c.ProxyConnection.Get(), "NEWNYM")
		},
	}
	c.StatusCommand = dependency.LazyDependency[*commands.StatusCommand]{
		InitFunc: func() *commands.StatusCommand {
			factory := c.HttpFactory.Get()
			proxyHost := c.Config.Get().Proxy.Host
			proxyPort := c.Config.Get().Proxy.Port
			timeout := time.Duration(10) * time.Second
			return commands.NewStatusCommand(proxyHost, proxyPort, factory, timeout)
		},
	}

	// Proxy services
	c.ProxyService = dependency.LazyDependency[*services.Service]{
		InitFunc: func() *services.Service {
			factory := c.HttpFactory.Get()
			proxyHost := c.Config.Get().Proxy.Host
			proxyPort := c.Config.Get().Proxy.Port
			timeout := 10 * time.Second
			return services.NewService(factory, proxyHost, proxyPort, timeout)
		},
	}
	c.IdentityService = dependency.LazyDependency[*services.Identity]{
		InitFunc: func() *services.Identity {
			conn := c.ProxyConnection.Get()
			auth := c.AuthenticateCommand.Get()
			signal := c.SignalCommand.Get()
			status := c.StatusCommand.Get()
			strategy := c.RetryStrategy.Get()
			url := c.Config.Get().Proxy.PingUrl
			return services.NewIdentity(conn, auth, signal, status, strategy, url)
		},
	}
	c.CircuitManager = dependency.LazyDependency[*circuit.Manager]{
		InitFunc: func() *circuit.Manager {
			return circuit.NewManager(c.IdentityService.Get(), c.ProxyService.Get(), c.Config.Get().Proxy.PingUrl)
		},
	}

	return c
}
