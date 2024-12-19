package source

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockHandler is a mock implementation of the Handler interface for testing purposes.
type MockHandler struct{}

// ProcessURLs is a mock implementation of the ProcessURLs.
func (m *MockHandler) ProcessURLs(ctx context.Context) error { return nil }

// ProcessHTML is a mock implementation of the ProcessHTML.
func (m *MockHandler) ProcessHTML(ctx context.Context, batchSize int) error { return nil }

// TestFactory_RegisterAndRetrieveHandler tests the registration and retrieval of handlers from the Factory.
func TestFactory_RegisterAndRetrieveHandler(t *testing.T) {
	container := SetupTestContainer()
	factory := container.SourceFactory.Get()

	mockHandler := &MockHandler{}

	// Register a handler
	err := factory.Register("mockSource", mockHandler)
	require.NoError(t, err, "Failed to register handler")

	// Retrieve the handler
	handler, err := factory.Get("mockSource")
	require.NoError(t, err, "Failed to get handler")
	assert.Equal(t, mockHandler, handler, "Retrieved handler should match the registered handler")
}

// TestFactory_DuplicateRegistration tests that duplicate registrations return an error.
func TestFactory_DuplicateRegistration(t *testing.T) {
	container := SetupTestContainer()
	factory := container.SourceFactory.Get()

	mockHandler := &MockHandler{}

	// Register a handler
	err := factory.Register("mockSource", mockHandler)
	require.NoError(t, err, "Failed to register handler")

	// Attempt to register another handler with the same name
	err = factory.Register("mockSource", mockHandler)
	require.Error(t, err, "Expected an error when registering a duplicate handler")
	assert.Contains(t, err.Error(), "handler with name mockSource already exists", "Error should indicate duplicate registration")
}

// TestFactory_GetNonExistentHandler tests retrieval of a non-existent handler.
func TestFactory_GetNonExistentHandler(t *testing.T) {
	container := SetupTestContainer()
	factory := container.SourceFactory.Get()

	// Attempt to retrieve a non-existent handler
	handler, err := factory.Get("nonExistentHandler")
	require.Error(t, err, "Expected an error when getting non existent handler")
	assert.Nil(t, handler, "Handler should be nil")
	assert.Contains(t, err.Error(), "handler with name nonExistentHandler does not exist", "Error should indicate non-existent handler")
}

// TestFactory_GetAllHandlers tests retrieval of all registered handlers.
func TestFactory_GetAllHandlers(t *testing.T) {
	container := SetupTestContainer()
	factory := container.SourceFactory.Get()

	// Handlers
	source1 := &MockHandler{}
	source2 := &MockHandler{}

	// Register handlers
	err := factory.Register("source1", source1)
	require.NoError(t, err, "Failed to register handler 1")
	err = factory.Register("source2", source2)
	require.NoError(t, err, "Failed to register handler 2")

	// Retrieve all handlers
	handlers := factory.GetAllHandlers()
	assert.Len(t, handlers, 2, "Expected two handlers to be registered")
	assert.Equal(t, source1, handlers["source1"], "Handler 1 should match the registered handler")
	assert.Equal(t, source2, handlers["source2"], "Handler 2 should match the registered handler")
}
