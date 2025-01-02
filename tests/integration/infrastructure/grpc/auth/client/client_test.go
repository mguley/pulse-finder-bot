package client

import (
	"context"
	authv1 "infrastructure/proto/auth/gen"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestAuthClient_GenerateToken tests the GenerateToken method of the AuthClient.
//
// This test validates the following scenarios:
// 1. Valid request returns a valid token.
// 2. Missing issuer returns an InvalidArgument error.
// 3. Empty scopes return an InvalidArgument error.
func TestAuthClient_GenerateToken(t *testing.T) {
	// Set up the test environment with mock server and client.
	env := SetupTestEnvironment(t)

	// Define test cases.
	tests := []struct {
		name        string
		request     *authv1.GenerateTokenRequest
		expectedErr bool
		errCode     codes.Code
	}{
		{
			name: "Valid Request",
			request: &authv1.GenerateTokenRequest{
				Issuer: "test-issuer",
				Scopes: []string{"scope1", "scope2"},
			},
			expectedErr: false,
		},
		{
			name: "Missing Issuer",
			request: &authv1.GenerateTokenRequest{
				Scopes: []string{"scope1", "scope2"},
			},
			expectedErr: true,
			errCode:     codes.InvalidArgument,
		},
		{
			name: "Empty Scopes",
			request: &authv1.GenerateTokenRequest{
				Issuer: "test-issuer",
				Scopes: []string{},
			},
			expectedErr: true,
			errCode:     codes.InvalidArgument,
		},
	}

	// Execute the test cases.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			// Make the gRPC call.
			resp, err := env.Client.GenerateToken(ctx, tc.request.Issuer, tc.request.Scopes)

			if tc.expectedErr {
				// Validate the error.
				require.Error(t, err, "expected an error but got none")
				st, ok := status.FromError(err)
				require.True(t, ok, "error is not a gRPC status")
				assert.Equal(t, tc.errCode, st.Code(), "unexpected gRPC status code")
			} else {
				// Validate the response.
				require.NoError(t, err, "unexpected error")
				assert.NotNil(t, resp, "response should not be nil")
				assert.NotEmpty(t, resp, "token should not be empty")
			}
		})
	}
}
