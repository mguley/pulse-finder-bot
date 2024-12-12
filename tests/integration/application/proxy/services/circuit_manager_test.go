package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestManager_ChangeCircuit_Success tests the successful circuit change process.
func TestManager_ChangeCircuit_Success(t *testing.T) {
	container := SetupTestContainer(t)
	manager := container.CircuitManager.Get()

	// Request a new circuit and validate the result.
	ip, err := manager.ChangeCircuit()
	require.NoError(t, err, "ChangeCircuit should succeed without errors")
	assert.NotEmpty(t, ip, "should return a valid ip")
}

// TestManager_ChangeCircuit_VerificationFailure tests the behavior when circuit verification fails.
func TestManager_ChangeCircuit_VerificationFailure(t *testing.T) {
	container := SetupTestContainer(t)
	// Simulate a verification failure by setting an invalid ping URL.
	container.Config.Get().Proxy.PingUrl = "http://invalid-url"

	manager := container.CircuitManager.Get()

	// Request a new circuit with invalid ping URL.
	ip, err := manager.ChangeCircuit()
	require.Error(t, err, "ChangeCircuit should fail due to verification error")
	assert.Empty(t, ip, "should return an empty ip")
	assert.Contains(t, err.Error(), "identity request", "Error message should indicate verification failure")
}
