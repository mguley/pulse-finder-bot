package application

import (
	"application/config"
	"application/dependency"
	"infrastructure"
)

// Container is a struct that holds all the dependencies for the application.
// It acts as a central registry for services, ensuring that dependencies are managed in a lazy loaded manner.
type Container struct {
	Config                  dependency.LazyDependency[*config.Config]
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

	return c
}
