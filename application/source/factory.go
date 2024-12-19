package source

import (
	"domain/source"
	"fmt"
	"sync"
)

// Factory is responsible for managing and providing source handlers.
type Factory struct {
	mu       sync.RWMutex
	handlers map[string]source.Handler // Map of source names to their handlers.
}

// NewFactory creates and returns a new Factory instance.
func NewFactory() *Factory {
	return &Factory{
		handlers: make(map[string]source.Handler),
	}
}

// Register adds a new source handler to the factory.
func (f *Factory) Register(name string, handler source.Handler) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.handlers[name]; exists {
		return fmt.Errorf("handler with name %s already exists", name)
	}
	f.handlers[name] = handler
	return nil
}

// Get retrieves a source handler by name.
func (f *Factory) Get(name string) (source.Handler, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	handler, exists := f.handlers[name]
	if !exists {
		return nil, fmt.Errorf("handler with name %s does not exist", name)
	}
	return handler, nil
}

// GetAllHandlers retrieves all registered source handlers.
func (f *Factory) GetAllHandlers() map[string]source.Handler {
	f.mu.RLock()
	defer f.mu.RUnlock()

	handlers := make(map[string]source.Handler, len(f.handlers))
	for name, handler := range f.handlers {
		handlers[name] = handler
	}
	return handlers
}
