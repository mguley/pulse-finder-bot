package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// authenticate authenticates with a proxy control port using the AuthenticateCommand.
func authenticate(t *testing.T, c *TestContainer) {
	authCmd := c.AuthenticateCommand.Get()

	// Execute the command
	err := authCmd.Execute()
	require.NoError(t, err, "Authentication failed, cannot proceed with signal tests")
}

// TestSignalCommand_Success tests successfully sending a valid signal to the proxy.
func TestSignalCommand_Success(t *testing.T) {
	container := SetupTestContainer(t)
	authenticate(t, container)

	cmd := container.SignalCommand.Get()

	// Execute the command
	err := cmd.Execute()

	// Assert successful execution
	require.NoError(t, err, "Signal command failed unexpectedly")
}
