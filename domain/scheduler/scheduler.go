package scheduler

import "context"

// Scheduler defines the interface for a periodic task scheduler.
type Scheduler interface {
	// Start initiates the periodic execution of scheduled tasks.
	Start(ctx context.Context)

	// Stop gracefully terminates the scheduler, ensuring that ongoing tasks complete.
	Stop()
}
