package main

import (
	"application"
	"context"
	"fmt"
	"infrastructure/grpc/auth/client"
	vacancyClient "infrastructure/grpc/vacancy/client"
	vacancyv1 "infrastructure/proto/vacancy/gen"
	"log"
	"time"

	"google.golang.org/grpc/metadata"
)

const (
	DefaultTimeout = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	// Initialize the application's container.
	app := application.NewContainer()
	config := app.Config.Get()

	aClient := app.InfrastructureContainer.Get().AuthClient.Get()
	vClient := app.InfrastructureContainer.Get().VacancyClient.Get()
	defer cleanup(aClient, vClient)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	// Generate JWT token
	jwtToken, err := generateToken(ctx, aClient, config.AuthServer.Issuer, []string{"write"})
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}
	log.Printf("JWT token: %s", jwtToken)

	// Attach token to context
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", jwtToken))

	// Create vacancy
	v, err := createVacancy(
		ctx,
		vClient,
		"Software Engineer",
		"Tech Corp",
		"Develop and maintain software.",
		"2025-01-01",
		"remote")
	if err != nil {
		return fmt.Errorf("create vacancy: %w", err)
	}
	log.Printf("Vacancy: %v", v)
	return nil
}

func cleanup(authClient *client.AuthClient, vacancyClient *vacancyClient.VacancyClient) {
	var err error
	if err = authClient.Close(); err != nil {
		log.Printf("Error closing AuthClient: %v", err)
	}
	if err = vacancyClient.Close(); err != nil {
		log.Printf("Error closing VacancyClient: %v", err)
	}
}

func generateToken(
	ctx context.Context,
	authClient *client.AuthClient,
	issuer string,
	scopes []string,
) (string, error) {
	token, err := authClient.GenerateToken(ctx, issuer, scopes)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}

func createVacancy(
	ctx context.Context,
	vacancyClient *vacancyClient.VacancyClient,
	title, company, description, postedAt, location string,
) (*vacancyv1.CreateVacancyResponse, error) {
	v, err := vacancyClient.CreateVacancy(ctx, title, company, description, postedAt, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create vacancy: %w", err)
	}
	return v, nil
}
