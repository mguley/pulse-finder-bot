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

// Subcommand is an enum-like type for handling the different CLI subcommands.
type Subcommand int

const (
	SubcommandUnknown Subcommand = iota
	SubcommandCreate
	SubcommandDelete
	SubcommandPurge
)

// CreateOptions holds the CLI arguments for the `create` subcommand.
type CreateOptions struct {
	Title       string
	Company     string
	Description string
	PostedAt    string
	Location    string
}

// DeleteOptions holds the CLI arguments for the `delete` subcommand.
type DeleteOptions struct {
	ID int64
}

func main() {
	// Parse subcommand.
	subcommand, createOpts, deleteOpts, err := parseSubcommand(os.Args)
	if err != nil {
		log.Println(err)
		printUsage()
		os.Exit(1)
	}
	if subcommand == SubcommandUnknown {
		printUsage()
		os.Exit(1)
	}

	// Initialize application container and clients.
	app := application.NewContainer()
	cfg := app.Config.Get()

	aClient := app.InfrastructureContainer.Get().AuthClient.Get()
	vClient := app.InfrastructureContainer.Get().VacancyClient.Get()
	defer cleanup(aClient, vClient)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	// Acquire JWT token
	jwtToken, err := generateToken(ctx, aClient, cfg.AuthServer.Issuer, []string{"write"})
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return
	}
	log.Printf("JWT: %s", jwtToken)

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", jwtToken))
	if err = handleSubcommand(ctx, subcommand, vClient, createOpts, deleteOpts); err != nil {
		log.Printf("Error running subcommand: %v", err)
		return
	}
}

// parseSubcommand inspects os.Args, determines which subcommand is being used, and parses any associated flags.
// Returns a Subcommand type and any subcommand-specific options structs.
func parseSubcommand(args []string) (Subcommand, *CreateOptions, *DeleteOptions, error) {
	if len(args) < 2 {
		// No subcommand provided
		return SubcommandUnknown, nil, nil, errors.New("no subcommand provided")
	}
	subcommandStr := args[1]

	// Define flag sets for each subcommand
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	purgeCmd := flag.NewFlagSet("purge", flag.ExitOnError)

	// Create subcommand flags
	createOpts := &CreateOptions{}
	createCmd.StringVar(&createOpts.Title, "title", "", "Job vacancy title")
	createCmd.StringVar(&createOpts.Company, "company", "", "Company name")
	createCmd.StringVar(&createOpts.Description, "desc", "", "Job vacancy description")
	createCmd.StringVar(&createOpts.PostedAt, "postedAt", "", "Posted date (format: YYYY-MM-DD)")
	createCmd.StringVar(&createOpts.Location, "location", "", "Vacancy location")

	// Delete subcommand flags
	deleteOpts := &DeleteOptions{}
	deleteCmd.Int64Var(&deleteOpts.ID, "id", 0, "Vacancy ID to delete")

	// Determine which subcommand and parse flags
	switch subcommandStr {
	case "create":
		if err := createCmd.Parse(args[2:]); err != nil {
			return SubcommandUnknown, nil, nil, fmt.Errorf("error parsing create flags: %w", err)
		}
		return SubcommandCreate, createOpts, nil, nil

	case "delete":
		if err := deleteCmd.Parse(args[2:]); err != nil {
			return SubcommandUnknown, nil, nil, fmt.Errorf("error parsing delete flags: %w", err)
		}
		return SubcommandDelete, nil, deleteOpts, nil

	case "purge":
		if err := purgeCmd.Parse(args[2:]); err != nil {
			return SubcommandUnknown, nil, nil, fmt.Errorf("error parsing purge flags: %w", err)
		}
		return SubcommandPurge, nil, nil, nil

	default:
		return SubcommandUnknown, nil, nil, fmt.Errorf("unknown subcommand: %q", subcommandStr)
	}
}

// handleSubcommand dispatches the correct subcommand execution based on the Subcommand enum.
func handleSubcommand(
	ctx context.Context,
	subcommand Subcommand,
	vClient *vacancyClient.VacancyClient,
	createOpts *CreateOptions,
	deleteOpts *DeleteOptions,
) error {
	switch subcommand {
	case SubcommandCreate:
		return runCreate(ctx, vClient, createOpts)
	case SubcommandDelete:
		return runDelete(ctx, vClient, deleteOpts)
	case SubcommandPurge:
		return runPurge(ctx, vClient)
	default:
		return errors.New("unsupported subcommand")
	}
}

// runCreate contains the logic for the "create" subcommand.
func runCreate(
	ctx context.Context,
	vClient *vacancyClient.VacancyClient,
	opts *CreateOptions,
) error {
	if opts.Title == "" || opts.Company == "" {
		return errors.New("both --title and --company are required for create")
	}

	resp, err := vClient.CreateVacancy(ctx, opts.Title, opts.Company, opts.Description, opts.PostedAt, opts.Location)
	if err != nil {
		return fmt.Errorf("create vacancy: %w", err)
	}
	log.Printf("Vacancy created successfully: %+v", resp)
	return nil
}

// runDelete contains the logic for the "delete" subcommand.
func runDelete(ctx context.Context, vClient *vacancyClient.VacancyClient, opts *DeleteOptions) error {
	if opts.ID <= 0 {
		return errors.New("vacancy --id must be a positive integer")
	}

	resp, err := vClient.DeleteVacancy(ctx, opts.ID)
	if err != nil {
		return fmt.Errorf("delete vacancy: %w", err)
	}
	log.Printf("Vacancy deleted successfully. Response: %+v", resp)
	return nil
}

// runPurge contains the logic for the "purge" subcommand.
func runPurge(ctx context.Context, vClient *vacancyClient.VacancyClient) error {
	resp, err := vClient.PurgeVacancies(ctx)
	if err != nil {
		return fmt.Errorf("purge vacancies: %w", err)
	}
	log.Printf("All vacancies purged successfully. Response: %+v", resp)
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
  purge    Remove all job vacancies

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

  # Purge all vacancies:
  %[1]s purge

`, os.Args[0])
}
