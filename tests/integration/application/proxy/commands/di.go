package commands

import (
	"application/config"
	"application/dependency"
	"application/proxy/commands"
	"infrastructure/proxy/port"
	"time"
)

// TestContainer holds dependencies for the integration tests.
type TestContainer struct {
	Config              dependency.LazyDependency[*config.Config]
	ProxyConnection     dependency.LazyDependency[*port.Connection]
	AuthenticateCommand dependency.LazyDependency[*commands.AuthenticateCommand]
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

	// Proxy commands
	c.AuthenticateCommand = dependency.LazyDependency[*commands.AuthenticateCommand]{
		InitFunc: func() *commands.AuthenticateCommand {
			return commands.NewAuthenticateCommand(c.ProxyConnection.Get())
		},
	}

	return c
}
