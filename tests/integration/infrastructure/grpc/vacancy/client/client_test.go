package client

import (
	"context"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestVacancyClient_CreateVacancy tests the CreateVacancy method of the VacancyClient.
//
// This test validates the following scenarios:
// 1. Valid request returns a successful response with the correct details.
// 2. Missing Title returns an InvalidArgument error.
// 3. Missing Company returns an InvalidArgument error.
// 4. Invalid PostedAt date format returns an InvalidArgument error.
func TestVacancyClient_CreateVacancy(t *testing.T) {
	// Set up the test environment with mock server and client.
	env := SetupTestEnvironment(t)

	// Define test cases.
	tests := []struct {
		name        string
		request     *vacancyv1.CreateVacancyRequest
		expectedErr bool
		errCode     codes.Code
	}{
		{
			name: "Valid Request",
			request: &vacancyv1.CreateVacancyRequest{
				Title:       "Software Engineer",
				Company:     "Tech Corp",
				Description: "Develop and maintain software.",
				PostedAt:    "2025-01-01",
				Location:    "Remote",
			},
			expectedErr: false,
		},
		{
			name: "Missing Title",
			request: &vacancyv1.CreateVacancyRequest{
				Company:     "Tech Corp",
				Description: "Develop and maintain software.",
				PostedAt:    "2025-01-01",
				Location:    "Remote",
			},
			expectedErr: true,
			errCode:     codes.InvalidArgument,
		},
		{
			name: "Missing Company",
			request: &vacancyv1.CreateVacancyRequest{
				Title:       "Software Engineer",
				Description: "Develop and maintain software.",
				PostedAt:    "2025-01-01",
				Location:    "Remote",
			},
			expectedErr: true,
			errCode:     codes.InvalidArgument,
		},
		{
			name: "Invalid PostedAt Date Format",
			request: &vacancyv1.CreateVacancyRequest{
				Title:       "Software Engineer",
				Company:     "Tech Corp",
				Description: "Develop and maintain software.",
				PostedAt:    "01-01-2025", // Incorrect format
				Location:    "Remote",
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
			resp, err := env.Client.CreateVacancy(
				ctx,
				tc.request.GetTitle(),
				tc.request.GetCompany(),
				tc.request.GetDescription(),
				tc.request.GetPostedAt(),
				tc.request.GetLocation())

			if tc.expectedErr {
				// Validate the error.
				require.Error(t, err, "expected error but got none")
				st, ok := status.FromError(err)
				require.True(t, ok, "error is not a gRPC status")
				assert.Equal(t, tc.errCode, st.Code(), "unexpected gRPC status code")
			} else {
				// Validate the response.
				require.NoError(t, err, "unexpected error")
				assert.NotNil(t, resp, "response should not be nil")
				assert.NotNil(t, resp.GetId(), "id should not be nil")
				assert.Equal(t, tc.request.GetTitle(), resp.GetTitle(), "title mismatch")
			}
		})
	}
}
