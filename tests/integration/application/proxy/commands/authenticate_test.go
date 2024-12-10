package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthenticateCommand_Success tests successful proxy authentication.
func TestAuthenticateCommand_Success(t *testing.T) {
	container := SetupTestContainer(t)
	cmd := container.AuthenticateCommand.Get()

	// Execute the command
	err := cmd.Execute()

	// Assert successful execution
	require.NoError(t, err, "Authentication command failed unexpectedly")
}

// TestAuthenticateCommand_InvalidPassword tests authentication failure due to an invalid password.
func TestAuthenticateCommand_InvalidPassword(t *testing.T) {
	container := SetupTestContainer(t)

	// Override the proxy password with an invalid one for the test
	container.Config.Get().Proxy.ControlPassword = "invalid_password"
	cmd := container.AuthenticateCommand.Get()

	// Execute the command
	err := cmd.Execute()

	// Assert failure
	require.Error(t, err, "Expected an authentication failure")
	assert.Contains(t, err.Error(), "authentication failed", "Expected error to indicate authentication failure")
}
