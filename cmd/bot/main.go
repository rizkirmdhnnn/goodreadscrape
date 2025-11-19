package main

import (
	"goodreadscrape/internal/app"
	"goodreadscrape/internal/config"
	"log"
	"os"
)

func main() {
	// Parse configuration
	cfg := config.ParseFlags()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Starting GoodScrape with config: APIKey=%s, Concurrency=%d, MaxReviews=%d, Language=%s, OutputFile=%s",
		cfg.APIKey, cfg.Concurrency, cfg.MaxReviews, cfg.Language, cfg.OutputFile)

	// Initialize and run the application
	application := app.NewScraperApp(cfg)
	application.Run()

	os.Exit(0)
}
