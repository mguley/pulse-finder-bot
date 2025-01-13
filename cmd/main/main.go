package main

import (
	"application"
	"context"
	"domain/source"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// HandlerRegistration holds the metadata for a source handler registration.
type HandlerRegistration struct {
	Name    string         // Unique name identifying the handler.
	Handler source.Handler // The actual handler instance to be registered.
}

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
	// Get the processor service and source factory from the application container.
	processor := c.ProcessorService.Get()
	sourceFactory := c.SourceFactory.Get()

	// Dynamically register each handler.
	for _, handler := range getItemsToProcess(c) {
		if err := sourceFactory.Register(handler.Name, handler.Handler); err != nil {
			return fmt.Errorf("register %s error: %v", handler.Name, err)
		}
		log.Printf("Registered source handler: %s", handler.Name)
	}

	// Start the processor.
	log.Println("Starting the processor...")
	processor.Run(ctx)
	log.Println("Processor completed successfully.")
	return nil
}

// getItemsToProcess defines the handlers to be registered and processed.
func getItemsToProcess(c *application.Container) []HandlerRegistration {
	return []HandlerRegistration{
		{
			Name:    "alfa",
			Handler: c.AlfaHandler.Get(),
		},
		{
			Name:    "beta",
			Handler: c.BetaHandler.Get(),
		},
	}
}

// main is the entry point for the application.
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
