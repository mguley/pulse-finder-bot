package server

import (
	"context"
	"fmt"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockVacancyService is a mock implementation of VacancyServiceServer for testing purposes.
// It implements the VacancyServiceServer interface and provides simulated behavior for the CreateVacancy method.
type MockVacancyService struct {
	vacancyv1.UnimplementedVacancyServiceServer // Ensures forward compatibility with the gRPC interface.
}

// NewMockVacancyService creates and returns a new instance of MockVacancyService.
func NewMockVacancyService() *MockVacancyService { return &MockVacancyService{} }

// CreateVacancy simulates the behavior of the CreateVacancy RPC method.
// It validates the incoming request and returns either a simulated response or an error.
// If the context is canceled, it returns a gRPC Canceled status.
func (s *MockVacancyService) CreateVacancy(ctx context.Context, req *vacancyv1.CreateVacancyRequest) (*vacancyv1.CreateVacancyResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context canceled")
	default:
		return s.handleCreateVacancy(ctx, req)
	}
}

// handleCreateVacancy validates the request and returns a mock CreateVacancyResponse or an error.
func (s *MockVacancyService) handleCreateVacancy(ctx context.Context, req *vacancyv1.CreateVacancyRequest) (*vacancyv1.CreateVacancyResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	// Return a simulated response.
	return &vacancyv1.CreateVacancyResponse{
		Id:          1,
		Title:       req.GetTitle(),
		Company:     req.GetCompany(),
		Description: req.GetDescription(),
		PostedAt:    req.GetPostedAt(),
		Location:    req.GetLocation(),
	}, nil
}

// validateRequest performs validation on the CreateVacancyRequest.
func (s *MockVacancyService) validateRequest(req *vacancyv1.CreateVacancyRequest) error {
	var validationErrors []error

	if err := s.validateStringField(req.Title, "title"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := s.validateStringField(req.Company, "company"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := s.validateDateField(req.PostedAt, "posted_at", "2006-01-02"); err != nil {
		validationErrors = append(validationErrors, err)
	}

	return s.combineErrors(validationErrors)
}

// validateStringField checks if a string field is provided and non-empty.
func (s *MockVacancyService) validateStringField(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return status.Errorf(codes.InvalidArgument,
			"%s must be provided and cannot be empty or whitespace", fieldName)
	}
	return nil
}

// validateDateField checks if a date field is in the expected format.
func (s *MockVacancyService) validateDateField(value, fieldName, format string) error {
	if strings.TrimSpace(value) == "" {
		return status.Errorf(codes.InvalidArgument,
			"%s must be provided and cannot be empty or whitespace", fieldName)
	}
	if _, err := time.Parse(format, value); err != nil {
		return status.Errorf(codes.InvalidArgument,
			"%s must be in the format %s. Example: %s",
			fieldName, format, time.Now().Format(format))
	}
	return nil
}

// combineErrors combines multiple errors into a single gRPC error.
func (s *MockVacancyService) combineErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var sb strings.Builder
	for i, err := range errs {
		sb.WriteString(fmt.Sprintf("Error %d: %v\n", i+1, err.Error()))
	}
	return status.Error(codes.InvalidArgument, sb.String())
}
