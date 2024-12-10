package infrastructure

import (
	"application/config"
	"application/dependency"
	"fmt"
	"infrastructure/proxy/port"
	"time"
)

// Container provides a lazily initialized set of dependencies for the infrastructure layer.
type Container struct {
	ProxyConnection dependency.LazyDependency[*port.Connection]
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

	return c
}
