package app

import (
	"context"
	"fmt"
	"goodreadscrape/internal/config"
	"goodreadscrape/internal/models"
	"goodreadscrape/internal/scraper"
	"goodreadscrape/internal/storage"
	"goodreadscrape/internal/validator"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// ScraperApp holds the application dependencies
type ScraperApp struct {
	Config    *config.Config
	Scraper   scraper.GoodreadsScraper
	Storage   storage.Storage
	saveMutex sync.Mutex
}

// NewScraperApp creates a new ScraperApp instance
func NewScraperApp(cfg *config.Config) *ScraperApp {
	return &ScraperApp{
		Config:  cfg,
		Scraper: scraper.NewGoodreadsScraper(cfg.APIKey, cfg.Verbose),
		Storage: storage.NewCSVStorage(),
	}
}

// Run executes the scraping process
func (app *ScraperApp) Run() {
	var urls []string
	var err error

	// Determine source of URLs
	if app.Config.InputURL != "" {
		urls = []string{app.Config.InputURL}
		log.Printf("Processing single URL: %s", app.Config.InputURL)
	} else if app.Config.InputFile != "" {
		urls, err = config.LoadURLsFromFile(app.Config.InputFile)
		if err != nil {
			log.Fatalf("Failed to load URLs from file: %v", err)
		}
		log.Printf("Loaded %d URLs from %s", len(urls), app.Config.InputFile)
	} else {
		log.Fatal("Error: You must provide either a single URL as an argument or an input file with -f")
	}

	// Validate URLs
	var validURLs []string
	for _, url := range urls {
		if validator.ValidateGoodreadsURL(url) {
			validURLs = append(validURLs, url)
		} else {
			log.Printf("Warning: Invalid Goodreads URL: %s", url)
		}
	}

	if len(validURLs) == 0 {
		log.Println("No valid URLs to process.")
		return
	}

	log.Printf("Processing %d valid URLs with %d workers...", len(validURLs), app.Config.Concurrency)

	// Create context that cancels on interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start worker pool
	jobs := make(chan string, len(validURLs))
	results := make(chan models.BookData, len(validURLs))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < app.Config.Concurrency; i++ {
		wg.Add(1)
		go app.worker(ctx, i+1, jobs, results, &wg)
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for _, url := range validURLs {
			select {
			case jobs <- url:
			case <-ctx.Done():
				log.Println("Signal received. Stopping new job dispatch...")
				return
			}
		}
	}()

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results
	successCount := 0
	processedCount := 0
	totalURLs := len(validURLs)

	for bookData := range results {
		processedCount++

		if len(bookData.Reviews) > 0 {
			app.saveMutex.Lock()
			err := app.Storage.SaveReviews(bookData.Reviews, app.Config.OutputFile)
			app.saveMutex.Unlock()

			if err != nil {
				log.Printf("❌ [%d/%d] Failed to save reviews for %s: %v", processedCount, totalURLs, bookData.Metadata.Title, err)
			} else {
				log.Printf("✅ [%d/%d] Saved %d reviews for '%s'", processedCount, totalURLs, len(bookData.Reviews), bookData.Metadata.Title)
				successCount++
			}
		} else {
			log.Printf("⚠️ [%d/%d] No reviews found for '%s'", processedCount, totalURLs, bookData.Metadata.Title)
			successCount++ // Count as success even if no reviews? Yes, scraping succeeded.
		}
	}

	fmt.Println("---------------------------------------------------------")
	log.Printf("🎉 Scraping completed! Successfully processed %d/%d URLs.", successCount, totalURLs)
	log.Printf("📂 Results saved to: %s", app.Config.OutputFile)
	fmt.Println("---------------------------------------------------------")
}

func (app *ScraperApp) worker(ctx context.Context, id int, jobs <-chan string, results chan<- models.BookData, wg *sync.WaitGroup) {
	defer wg.Done()

	filters := models.Filters{Language: app.Config.Language}

	for url := range jobs {
		// Check context before starting work (optional, as channel close handles it, but good for fast exit)
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Verbose logging
		if app.Config.Verbose {
			log.Printf("Worker %d: Starting scraping for %s", id, url)
		}

		// Scrape book data
		bookData, err := app.Scraper.ScrapeBookData(url, app.Config.MaxReviews, filters)
		if err != nil {
			log.Printf("❌ Worker %d: Failed to scrape %s: %v", id, url, err)
			continue
		}

		// Verbose logging
		if app.Config.Verbose {
			log.Printf("Worker %d: Finished scraping %s (%d reviews)", id, url, len(bookData.Reviews))
		}

		results <- bookData
	}
}
