package main

import (
	"application"
	"context"
	"errors"
	"flag"
	"fmt"
	"infrastructure/grpc/auth/client"
	vacancyClient "infrastructure/grpc/vacancy/client"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/metadata"
)

const (
	DefaultTimeout = 10 * time.Second
)

func main() {
	// -------------------------------------------------
	// Define subcommands
	// -------------------------------------------------
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// -------------------------------------------------
	// Define flags for the "create" subcommand
	// -------------------------------------------------
	createTitle := createCmd.String("title", "", "Job vacancy title (required)")
	createCompany := createCmd.String("company", "", "Company name (required)")
	createDescription := createCmd.String("desc", "", "Job vacancy description")
	createPostedAt := createCmd.String("postedAt", "", "Posted date (format: YYYY-MM-DD)")
	createLocation := createCmd.String("location", "", "Vacancy location")

	// -------------------------------------------------
	// Define flags for the "delete" subcommand
	// -------------------------------------------------
	deleteID := deleteCmd.Int64("id", 0, "Vacancy ID to delete (required)")

	// -------------------------------------------------
	// Check we have at least one subcommand
	// -------------------------------------------------
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// The subcommand is in os.Args[1] (e.g. "create" or "delete")
	subcommand := os.Args[1]

	// -------------------------------------------------
	// Parse subcommand flags
	// -------------------------------------------------
	var err error
	switch subcommand {
	case "create":
		err = createCmd.Parse(os.Args[2:])
	case "delete":
		err = deleteCmd.Parse(os.Args[2:])
	default:
		fmt.Printf("Unknown subcommand: %q\n\n", subcommand)
		printUsage()
		os.Exit(1)
	}
	if err != nil {
		log.Fatalf("Error parsing %q subcommand flags: %v", subcommand, err)
	}

	// -------------------------------------------------
	// Prepare environment and clients
	// -------------------------------------------------
	app := application.NewContainer()
	config := app.Config.Get()

	aClient := app.InfrastructureContainer.Get().AuthClient.Get()
	vClient := app.InfrastructureContainer.Get().VacancyClient.Get()
	defer cleanup(aClient, vClient)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	// -------------------------------------------------
	// Acquire JWT token
	// -------------------------------------------------
	jwtToken, err := generateToken(ctx, aClient, config.AuthServer.Issuer, []string{"write"})
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return
	}
	log.Printf("JWT: %s", jwtToken)

	// Attach the token to the outgoing context
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", jwtToken))

	// -------------------------------------------------
	// Execute the requested subcommand logic
	// -------------------------------------------------
	switch subcommand {
	case "create":
		err = runCreate(ctx, vClient, *createTitle, *createCompany, *createDescription, *createPostedAt, *createLocation)
	case "delete":
		err = runDelete(ctx, vClient, *deleteID)
	}

	if err != nil {
		log.Printf("Error running %q subcommand: %v", subcommand, err)
		return
	}
}

// runCreate contains the logic for the "create" subcommand.
func runCreate(
	ctx context.Context,
	vClient *vacancyClient.VacancyClient,
	title, company, desc, postedAt, location string,
) error {
	if title == "" || company == "" {
		return errors.New("both --title and --company are required for create")
	}

	resp, err := vClient.CreateVacancy(ctx, title, company, desc, postedAt, location)
	if err != nil {
		return fmt.Errorf("create vacancy: %w", err)
	}
	log.Printf("Vacancy created successfully: %+v", resp)
	return nil
}

// runDelete contains the logic for the "delete" subcommand.
func runDelete(ctx context.Context, vClient *vacancyClient.VacancyClient, id int64) error {
	if id <= 0 {
		return errors.New("vacancy --id must be a positive integer")
	}

	resp, err := vClient.DeleteVacancy(ctx, id)
	if err != nil {
		return fmt.Errorf("delete vacancy: %w", err)
	}
	log.Printf("Vacancy deleted successfully. Response: %+v", resp)
	return nil
}

// cleanup closes both the AuthClient and VacancyClient, logging any errors.
func cleanup(authClient *client.AuthClient, vacancyClient *vacancyClient.VacancyClient) {
	if err := authClient.Close(); err != nil {
		log.Printf("Error closing AuthClient: %v", err)
	}
	if err := vacancyClient.Close(); err != nil {
		log.Printf("Error closing VacancyClient: %v", err)
	}
}

// generateToken fetches a JWT token from the AuthClient.
func generateToken(ctx context.Context, authClient *client.AuthClient, issuer string, scopes []string) (string, error) {
	token, err := authClient.GenerateToken(ctx, issuer, scopes)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

// printUsage shows the usage for the entire CLI, including subcommand usage.
func printUsage() {
	fmt.Printf(`Usage:
  %s <subcommand> [options]

Subcommands:
  create   Create a job vacancy
  delete   Delete an existing job vacancy

Examples:
  # Create a new vacancy:
  %[1]s create \
      --title="Software Engineer" \
      --company="Tech Corp" \
      --desc="Develop and maintain software." \
      --postedAt="2025-01-01" \
      --location="Remote"

  # Delete a vacancy by ID:
  %[1]s delete --id=123

`, os.Args[0])
}
