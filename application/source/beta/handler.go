package beta

import (
	"application/proxy/circuit"
	"application/url/sitemap"
	"context"
	"domain/html"
	"domain/url/entity"
	urlRepository "domain/url/repository"
	vacancyEntity "domain/vacancy/entity"
	vacancyRepository "domain/vacancy/repository"
	"fmt"
	"time"
)

// Handler processes URLs and HTML content.
type Handler struct {
	url               string                              // Base URL of the sitemap to process.
	sitemapService    *sitemap.Service                    // Service processes the URLs from a sitemap.
	circuitManager    *circuit.Manager                    // Service manages the proxy circuit lifecycle.
	urlRepository     urlRepository.UrlRepository         // Service manages URL entities in the data source.
	vacancyRepository vacancyRepository.VacancyRepository // Service handles the storage of parsed vacancy details.
	fetcher           html.Fetcher                        // Service fetches HTML content over HTTP.
	parser            html.Parser                         // Service extracts vacancy details from raw HTML content.
}

// NewHandler creates and returns a new Handler instance.
func NewHandler(
	url string,
	sitemapService *sitemap.Service,
	circuitManager *circuit.Manager,
	urlRepo urlRepository.UrlRepository,
	vacancyRepo vacancyRepository.VacancyRepository,
	fetcher html.Fetcher,
	parser html.Parser,
) *Handler {
	return &Handler{
		url:               url,
		sitemapService:    sitemapService,
		circuitManager:    circuitManager,
		urlRepository:     urlRepo,
		vacancyRepository: vacancyRepo,
		fetcher:           fetcher,
		parser:            parser,
	}
}

// ProcessURLs retrieves and processes sitemap URLs.
func (h *Handler) ProcessURLs(ctx context.Context) error {
	if err := h.sitemapService.ProcessUrls(ctx, h.url); err != nil {
		return fmt.Errorf("process urls: %w", err)
	}
	return nil
}

// ProcessHTML processes URLs in batches with a delay.
func (h *Handler) ProcessHTML(ctx context.Context, batchSize int) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Println("Sleeping for 15 seconds...")
			time.Sleep(15 * time.Second)

			// Process a batch and determine if there are more URLs to process
			hasMore, err := h.processBatch(ctx, batchSize)
			if err != nil {
				return fmt.Errorf("process batch: %w", err)
			}

			// Exit if there are no more URLs to process
			if !hasMore {
				fmt.Println("All URLs processed successfully")
				return nil
			}
		}
	}
}

// processBatch retrieves and processes a batch of URLs.
func (h *Handler) processBatch(ctx context.Context, batchSize int) (bool, error) {
	// Fetch a batch of unprocessed URLs with the "pending" status.
	urls, err := h.urlRepository.FetchBatch(ctx, "pending", batchSize)
	if err != nil {
		return false, fmt.Errorf("fetch batch: %w", err)
	}

	// Exit if there are no URLs to process.
	if len(urls) == 0 {
		fmt.Println("No more URLs to process")
		return false, nil
	}

	// Switch proxy circuit
	if err = h.switchCircuit(); err != nil {
		return true, fmt.Errorf("switch circuit: %w", err)
	}

	// Process each URL in the batch
	for _, url := range urls {
		if err = h.processSingleURL(ctx, url); err != nil {
			fmt.Printf("process url %s failed: %v", url, err)
		}
	}
	return true, nil
}

// processSingleURL fetches, parses, and saves data for a single URL.
func (h *Handler) processSingleURL(ctx context.Context, url *entity.Url) error {
	processedTime := time.Now()

	// Fetch raw HTML content from the URL.
	raw, err := h.fetcher.Fetch(ctx, url.Address)
	if err != nil {
		if err = h.markStatus(ctx, url, "failed", processedTime); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	// Parse the fetched HTML into structured format.
	data, err := h.parser.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse url, %s: %w", url.Address, err)
	}
	defer data.Release()

	// Map DTO to entity and save
	v := &vacancyEntity.Vacancy{}
	data.ToEntity(v)
	if err = h.vacancyRepository.Save(ctx, v); err != nil {
		return fmt.Errorf("save vacancy, %s: %w", v, err)
	}

	// Update URL status
	if err = h.markStatus(ctx, url, "success", processedTime); err != nil {
		return fmt.Errorf("mark vacancy, %s: %w", url, err)
	}

	fmt.Printf("Processed URL %s\n", url.Address)
	return nil
}

// switchCircuit switches the proxy circuit and logs the new IP.
func (h *Handler) switchCircuit() error {
	ip, err := h.circuitManager.ChangeCircuit()
	if err != nil {
		return fmt.Errorf("change circuit: %w", err)
	}
	fmt.Printf("Using new circuit ip: %s\n", ip)
	return nil
}

// markStatus marks the processed URL with a provided status.
func (h *Handler) markStatus(ctx context.Context, url *entity.Url, status string, processedTime time.Time) error {
	if err := h.urlRepository.UpdateStatus(ctx, url.ID.Hex(), status, &processedTime); err != nil {
		return fmt.Errorf("update url status: %w", err)
	}
	return nil
}
