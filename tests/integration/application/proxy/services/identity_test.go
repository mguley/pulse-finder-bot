package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIdentity_Request validates the end-to-end functionality of the Identity service's Request method.
func TestIdentity_Request(t *testing.T) {
	container := SetupTestContainer()
	identity := container.IdentityService.Get()

	// Request new identity
	err := identity.Request()
	assert.NoError(t, err, "Identity request should succeed")
}
