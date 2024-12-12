package application

import (
	"application/config"
	"application/dependency"
	"application/proxy/circuit"
	"application/proxy/commands"
	"application/proxy/services"
	"application/proxy/strategies"
	"infrastructure"
	"time"
)

// Container is a struct that holds all the dependencies for the application.
// It acts as a central registry for services, ensuring that dependencies are managed in a lazy loaded manner.
type Container struct {
	Config                  dependency.LazyDependency[*config.Config]
	ProxyService            dependency.LazyDependency[*services.Service]
	RetryStrategy           dependency.LazyDependency[strategies.RetryStrategy]
	IdentityService         dependency.LazyDependency[*services.Identity]
	CircuitManager          dependency.LazyDependency[*circuit.Manager]
	AuthenticateCommand     dependency.LazyDependency[*commands.AuthenticateCommand]
	SignalCommand           dependency.LazyDependency[*commands.SignalCommand]
	StatusCommand           dependency.LazyDependency[*commands.StatusCommand]
	InfrastructureContainer dependency.LazyDependency[*infrastructure.Container]
}

// NewContainer creates and returns a new instance of Container.
// Each dependency is configured to initialize only when first accessed.
func NewContainer() *Container {
	c := &Container{}

	// Create container with base dependencies
	c.Config = dependency.LazyDependency[*config.Config]{
		InitFunc: config.LoadConfig,
	}

	// Domain/layer containers
	c.InfrastructureContainer = dependency.LazyDependency[*infrastructure.Container]{
		InitFunc: func() *infrastructure.Container {
			return infrastructure.NewContainer(c.Config.Get())
		},
	}

	// Proxy commands
	c.AuthenticateCommand = dependency.LazyDependency[*commands.AuthenticateCommand]{
		InitFunc: func() *commands.AuthenticateCommand {
			conn := c.InfrastructureContainer.Get().ProxyConnection.Get()
			return commands.NewAuthenticateCommand(conn)
		},
	}
	c.SignalCommand = dependency.LazyDependency[*commands.SignalCommand]{
		InitFunc: func() *commands.SignalCommand {
			conn := c.InfrastructureContainer.Get().ProxyConnection.Get()
			return commands.NewSignalCommand(conn, "NEWNYM")
		},
	}
	c.StatusCommand = dependency.LazyDependency[*commands.StatusCommand]{
		InitFunc: func() *commands.StatusCommand {
			factory := c.InfrastructureContainer.Get().HttpFactory.Get()
			client := c.InfrastructureContainer.Get().Socks5Client.Get()
			host := c.Config.Get().Proxy.Host
			port := c.Config.Get().Proxy.Port
			timeout := 10 * time.Second
			return commands.NewStatusCommand(host, port, client, factory, timeout)
		},
	}

	// Proxy services
	c.ProxyService = dependency.LazyDependency[*services.Service]{
		InitFunc: func() *services.Service {
			factory := c.InfrastructureContainer.Get().HttpFactory.Get()
			host := c.Config.Get().Proxy.Host
			port := c.Config.Get().Proxy.Port
			timeout := 10 * time.Second
			return services.NewService(factory, host, port, timeout)
		},
	}
	c.RetryStrategy = dependency.LazyDependency[strategies.RetryStrategy]{
		InitFunc: func() strategies.RetryStrategy {
			baseDelay := 5 * time.Second
			maxDelay := 30 * time.Second
			maxAttempts := 5
			multiplier := 2.0
			return strategies.NewExponentialBackoffStrategy(baseDelay, maxDelay, maxAttempts, multiplier)
		},
	}
	c.IdentityService = dependency.LazyDependency[*services.Identity]{
		InitFunc: func() *services.Identity {
			conn := c.InfrastructureContainer.Get().ProxyConnection.Get()
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
