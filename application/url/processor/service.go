package processor

import (
	"application/source"
	"context"
	"fmt"
)

// Service orchestrates the processing of URLs by managing multiple sources.
type Service struct {
	factory   *source.Factory // Factory for managing source handlers.
	batchSize int             // Batch size for processing.
}

// NewService creates and initializes a new Service instance.
func NewService(factory *source.Factory, batchSize int) *Service {
	return &Service{factory: factory, batchSize: batchSize}
}

// Run iterates through all registered sources and processes their URLs and HTML content.
func (s *Service) Run(ctx context.Context) {
	for {
		fmt.Println("Starting to process all registered sources.")
		sources := s.factory.GetAllHandlers()

		for name, handler := range sources {
			fmt.Printf("Processing source: %s\n", name)

			// Process URLs for the source.
			if err := handler.ProcessURLs(ctx); err != nil {
				fmt.Printf("Error processing URLs for source %s: %v\n", name, err)
				continue
			}

			// Process HTML for the source in batches.
			if err := handler.ProcessHTML(ctx, s.batchSize); err != nil {
				fmt.Printf("Error processing HTML for source %s: %v\n", name, err)
			}
		}

		// Done.
		fmt.Println("Finished processing all registered sources.")
		break
	}
}
