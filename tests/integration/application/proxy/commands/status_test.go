package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestStatusCommand_Success tests the successful execution of the StatusCommand.
func TestStatusCommand_Success(t *testing.T) {
	container := SetupTestContainer(t)
	cmd := container.StatusCommand.Get()

	// Execute the command
	res, err := cmd.Execute("https://httpbin.org/ip")

	// Assert successful execution
	require.NoError(t, err, "Status command failed unexpectedly")
	require.NotNil(t, res, "Status command returned nil")
	require.Contains(t, res, "origin", "Status command returned unexpected value")
}
