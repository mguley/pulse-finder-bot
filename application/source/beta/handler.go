package beta

import (
	"application/proxy/circuit"
	"application/url/processor/dto"
	"application/url/sitemap"
	"context"
	"domain/html"
	"domain/url/entity"
	urlRepository "domain/url/repository"
	vacancyEntity "domain/vacancy/entity"
	vacancyRepository "domain/vacancy/repository"
	"fmt"
	"sync"
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
func (h *Handler) ProcessURLs(ctx context.Context) (err error) {
	if err = h.sitemapService.ProcessUrls(ctx, h.url); err != nil {
		return fmt.Errorf("process urls: %w", err)
	}
	return nil
}

// ProcessHTML processes URLs in batches with a delay.
func (h *Handler) ProcessHTML(ctx context.Context, batchSize int) (err error) {
	var (
		maxConcurrency = 5
		hasMore        bool
		delayTime      = time.Duration(15) * time.Second
	)

	for {
		// Attempt to process a batch of URLs
		if hasMore, err = h.processBatch(ctx, batchSize, maxConcurrency); err != nil {
			return fmt.Errorf("process batch: %w", err)
		}
		if !hasMore {
			return nil
		}

		// Delay before the next iteration
		fmt.Println("Sleeping...")
		time.Sleep(delayTime)
	}
}

// processBatch fetches a batch of URLs and processes them respecting the maxConcurrency limit.
func (h *Handler) processBatch(ctx context.Context, batchSize, maxConcurrency int) (hasMore bool, err error) {
	var (
		urls         []*entity.Url
		status       = "pending"
		switchResult string
	)

	// Change identity
	if switchResult, err = h.circuitManager.ChangeCircuit(); err != nil {
		return false, fmt.Errorf("change circuit: %w", err)
	}
	fmt.Printf("Switch Result: %s\n", switchResult)

	if urls, err = h.urlRepository.FetchBatch(ctx, status, batchSize); err != nil {
		return false, fmt.Errorf("fetch batch: %w", err)
	}
	if len(urls) == 0 {
		return false, nil
	}

	var (
		semaphore = make(chan struct{}, maxConcurrency)
		wg        sync.WaitGroup
	)

	wg.Add(len(urls))
	for _, url := range urls {
		semaphore <- struct{}{}

		go func(entity *entity.Url) {
			defer wg.Done()
			defer func() { <-semaphore }()
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[ERROR] Recovered from panic in processUrl for %s: %v\n", entity.Address, r)
				}
			}()

			// Workload
			if pErr := h.processUrl(ctx, entity); pErr != nil {
				fmt.Printf("[WARN] failed to process URL %s: %v\n", entity.Address, pErr)
				return
			}
		}(url)
	}
	wg.Wait()

	// There might be more items in subsequent batches
	return true, nil
}

// processUrl fetches, parses, and saves data for a single URL.
func (h *Handler) processUrl(ctx context.Context, url *entity.Url) (err error) {
	var (
		processedTime = time.Now()
		body          string
		result        *dto.Vacancy
		vacancy       = &vacancyEntity.Vacancy{}
	)

	// Fetch raw HTML content from the URL.
	if body, err = h.fetcher.Fetch(ctx, url.Address); err != nil {
		if err = h.updateStatus(ctx, url, "failed", processedTime); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	// Parse the fetched HTML into structured format.
	if result, err = h.parser.Parse(body); err != nil {
		return fmt.Errorf("parse url, %s: %w", url.Address, err)
	}
	defer result.Release()

	// Save
	result.ToEntity(vacancy)
	if err = h.vacancyRepository.Save(ctx, vacancy); err != nil {
		return fmt.Errorf("save vacancy, %s: %w", vacancy, err)
	}

	// Update URL status
	if err = h.updateStatus(ctx, url, "success", processedTime); err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Printf("[INFO] successfully processed URL: %s\n", url.Address)
	return nil
}

// updateStatus updates information about the processed URL.
func (h *Handler) updateStatus(ctx context.Context, entity *entity.Url, status string, time time.Time) (err error) {
	if err = h.urlRepository.UpdateStatus(ctx, entity.ID.Hex(), status, &time); err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return nil
}
