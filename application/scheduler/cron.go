package scheduler

import (
	"context"
	"domain/vacancy/entity"
	"domain/vacancy/repository"
	"fmt"
	authClient "infrastructure/grpc/auth/client"
	vacancyClient "infrastructure/grpc/vacancy/client"
	"time"

	"google.golang.org/grpc/metadata"
)

// CronScheduler is a periodic job scheduler that transfers vacancies over gRPC.
type CronScheduler struct {
	repository    repository.VacancyRepository // Repository is used to fetch and update vacancy entities.
	authClient    *authClient.AuthClient       // Manages gRPC client connection to AuthService.
	vacancyClient *vacancyClient.VacancyClient // Manages gRPC client connection to VacancyService.
	done          chan struct{}                // Signal channel used to stop the scheduler’s loop.
	batchSize     int                          // Defines how many items are processed in a single batch.
	tokenIssuer   string                       // Is the issuer field for the JWT token generation.
	tokenScope    []string                     // Defines the scopes requested for the JWT token.
	tickerTime    time.Duration                // Interval at which the CronScheduler triggers its job.
}

// NewCronScheduler creates a new instance of CronScheduler.
func NewCronScheduler(
	repository repository.VacancyRepository,
	aClient *authClient.AuthClient,
	vClient *vacancyClient.VacancyClient,
	batchSize int,
	tokenIssuer string,
	tokenScope []string,
	tickerTime time.Duration,
) *CronScheduler {
	return &CronScheduler{
		repository:    repository,
		authClient:    aClient,
		vacancyClient: vClient,
		done:          make(chan struct{}),
		batchSize:     batchSize,
		tokenIssuer:   tokenIssuer,
		tokenScope:    tokenScope,
		tickerTime:    tickerTime,
	}
}

// Start initiates the periodic job using a time.Ticker.
func (s *CronScheduler) Start(ctx context.Context) {
	fmt.Println("starting cron scheduler...")
	ticker := time.NewTicker(s.tickerTime)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-s.done:
				fmt.Println("cron scheduler exit")
				return
			case <-ticker.C:
				err := s.transferVacancies(ctx)
				if err != nil {
					fmt.Printf("transfer vacancies err: %v\n", err)
				}
			}
		}
	}()
}

// Stop terminates the scheduler.
func (s *CronScheduler) Stop() {
	close(s.done)
}

// transferVacancies obtains a JWT token and processes the jobs in batches, sending them to the remote service.
func (s *CronScheduler) transferVacancies(ctx context.Context) error {
	// Generate a fresh JWT token for gRPC calls.
	jwtToken, err := s.generateToken(ctx)
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}

	// Attach the token to the outgoing gRPC metadata.
	outCtx := metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", jwtToken))

	// Keep processing batches until no more items remain.
	for {
		hasMore, err := s.processBatch(outCtx, s.batchSize)
		if err != nil {
			return fmt.Errorf("process batch: %w", err)
		}
		if !hasMore {
			fmt.Println("finished processing all batches")
			return nil
		}
	}
}

// processBatch retrieves and processes a batch of vacancies.
func (s *CronScheduler) processBatch(ctx context.Context, batchSize int) (bool, error) {
	items, err := s.repository.FetchBatch(ctx, batchSize)
	if err != nil {
		return false, fmt.Errorf("fetch batch: %w", err)
	}
	if len(items) == 0 {
		fmt.Println("no more items to process")
		return false, nil
	}

	for _, item := range items {
		if err = s.sendVacancy(ctx, item); err != nil {
			// Log error but do not return it so that the rest are still processed
			fmt.Printf("[WARN] could not send vacancy (ID=%s): %v\n", item.ID.Hex(), err)
		}
	}
	return true, nil
}

// sendVacancy calls the vacancyClient to create a vacancy via gRPC and updates the local entity's SentAt field.
func (s *CronScheduler) sendVacancy(ctx context.Context, item *entity.Vacancy) error {
	// Create the vacancy on the remote service via gRPC.
	_, err := s.vacancyClient.CreateVacancy(
		ctx,
		item.Title,
		item.Company,
		item.Description,
		item.PostedAt.Format(time.DateOnly),
		item.Location)
	if err != nil {
		return fmt.Errorf("send vacancy over gRPC: %w", err)
	}

	item.SentAt = time.Now()
	if err = s.repository.Update(ctx, item); err != nil {
		return fmt.Errorf("update SentAt field: %w", err)
	}
	return nil
}

// generateToken fetches a JWT token from the AuthClient using the configured issuer and scope.
func (s *CronScheduler) generateToken(ctx context.Context) (string, error) {
	token, err := s.authClient.GenerateToken(ctx, s.tokenIssuer, s.tokenScope)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}
