package commands

import (
	"application/config"
	"application/dependency"
	"application/proxy/commands"
	"domain/useragent"
	httpClient "infrastructure/http/client"
	"infrastructure/proxy/client"
	"infrastructure/proxy/client/agent"
	"infrastructure/proxy/port"
	"time"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config              dependency.LazyDependency[*config.Config]
	ProxyConnection     dependency.LazyDependency[*port.Connection]
	UserAgent           dependency.LazyDependency[useragent.Generator]
	Socks5Client        dependency.LazyDependency[*client.Socks5Client]
	HttpFactory         dependency.LazyDependency[*httpClient.Factory]
	AuthenticateCommand dependency.LazyDependency[*commands.AuthenticateCommand]
	SignalCommand       dependency.LazyDependency[*commands.SignalCommand]
	StatusCommand       dependency.LazyDependency[*commands.StatusCommand]
}

// NewTestContainer initializes a new test container.
func NewTestContainer() *TestContainer {
	c := &TestContainer{}

	c.Config = dependency.LazyDependency[*config.Config]{
		InitFunc: config.LoadConfig,
	}
	c.ProxyConnection = dependency.LazyDependency[*port.Connection]{
		InitFunc: func() *port.Connection {
			cfg := c.Config.Get()
			address := cfg.Proxy.Host + ":" + cfg.Proxy.ControlPort
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

	// Proxy commands
	c.AuthenticateCommand = dependency.LazyDependency[*commands.AuthenticateCommand]{
		InitFunc: func() *commands.AuthenticateCommand {
			return commands.NewAuthenticateCommand(c.ProxyConnection.Get())
		},
	}
	c.SignalCommand = dependency.LazyDependency[*commands.SignalCommand]{
		InitFunc: func() *commands.SignalCommand {
			return commands.NewSignalCommand(c.ProxyConnection.Get(), "NEWNYM")
		},
	}
	c.StatusCommand = dependency.LazyDependency[*commands.StatusCommand]{
		InitFunc: func() *commands.StatusCommand {
			f := c.HttpFactory.Get()
			s := c.Socks5Client.Get()
			h := c.Config.Get().Proxy.Host
			p := c.Config.Get().Proxy.Port
			t := 10 * time.Second
			return commands.NewStatusCommand(h, p, s, f, t)
		},
	}

	return c
}
