package client

import (
	"context"
	"errors"
	"fmt"
	vacancyv1 "infrastructure/proto/vacancy/gen"

	"google.golang.org/grpc"
)

// VacancyClient is a high-level wrapper over the underlying gRPC client connection to VacancyService.
// It manages the gRPC connection and provides methods to interact with the VacancyService.
type VacancyClient struct {
	conn   *grpc.ClientConn               // The underlying gRPC client connection.
	client vacancyv1.VacancyServiceClient // The generated VacancyService client.
}

// NewVacancyClient creates a new instance of VacancyClient based on the provided configuration.
func NewVacancyClient(env, address string) (*VacancyClient, error) {
	var (
		conn   *grpc.ClientConn
		config *Config
		err    error
	)

	switch env {
	case "prod":
		// Production environment with TLS (certificate file is optional here, adjust as needed).
		conn, config, err = NewGRPCClient(
			WithServerAddress(address),
			WithTLS(""))
	case "dev":
		// Development environment without TLS.
		conn, config, err = NewGRPCClient(
			WithServerAddress(address))
	default:
		return nil, errors.New("unsupported environment; must be \"prod\" or \"dev\"")
	}

	if err != nil {
		return nil, fmt.Errorf("create gRPC client: %w", err)
	}
	fmt.Printf("Connected to server: %s (TLS enabled: %v)\n", address, config.TLSEnabled)
	return &VacancyClient{
		conn:   conn,
		client: vacancyv1.NewVacancyServiceClient(conn),
	}, nil
}

// CreateVacancy calls the VacancyService's CreateVacancy RPC method.
// It returns the created vacancy details or an error if the call fails.
func (v *VacancyClient) CreateVacancy(
	ctx context.Context,
	title,
	company,
	description,
	postedAt,
	location string,
) (*vacancyv1.CreateVacancyResponse, error) {
	req := &vacancyv1.CreateVacancyRequest{
		Title:       title,
		Company:     company,
		Description: description,
		PostedAt:    postedAt,
		Location:    location,
	}
	resp, err := v.client.CreateVacancy(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create vacancy: %w", err)
	}

	return resp, nil
}

// DeleteVacancy calls the VacancyService's DeleteVacancy RPC method.
// It returns a confirmation message or an error if the call fails.
func (v *VacancyClient) DeleteVacancy(ctx context.Context, id int64) (*vacancyv1.DeleteVacancyResponse, error) {
	req := &vacancyv1.DeleteVacancyRequest{Id: id}

	resp, err := v.client.DeleteVacancy(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("delete vacancy: %w", err)
	}

	return resp, nil
}

// PurgeVacancies calls the VacancyService's PurgeVacancies RPC method.
// It removes all job vacancies from the database and returns a confirmation message or an error if the call fails.
func (v *VacancyClient) PurgeVacancies(ctx context.Context) (*vacancyv1.PurgeVacanciesResponse, error) {
	req := &vacancyv1.PurgeVacanciesRequest{}

	resp, err := v.client.PurgeVacancies(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("purge vacancies: %w", err)
	}

	return resp, nil
}

// Close closes the underlying gRPC connection.
func (v *VacancyClient) Close() error {
	return v.conn.Close()
}
