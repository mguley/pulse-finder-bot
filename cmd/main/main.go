package main

import (
	"application"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// setupGracefulShutdown handles termination signals.
func setupGracefulShutdown(cancelFunc context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-signalChan
		log.Println("Received termination signal, shutting down gracefully...")
		cancelFunc()
	}()
}

// runProcessor initializes and runs the application processor.
func runProcessor(ctx context.Context, c *application.Container) error {
	// Initialize the processor service and source factory.
	processor := c.ProcessorService.Get()
	sourceFactory := c.SourceFactory.Get()

	// Register the source handler.
	name := "beta"
	handler := c.BetaHandler.Get()

	if err := sourceFactory.Register(name, handler); err != nil {
		return fmt.Errorf("%w", err)
	}
	log.Printf("Registered source handler: %s", name)

	// Start the processor.
	log.Println("Starting the processor...")
	processor.Run(ctx)
	log.Println("Processor completed successfully.")
	return nil
}

func main() {
	// Initialize application container.
	c := application.NewContainer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown listener
	setupGracefulShutdown(cancel)

	if err := runProcessor(ctx, c); err != nil {
		log.Printf("Error running processor: %v", err)
	}
}
