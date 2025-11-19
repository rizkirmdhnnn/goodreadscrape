package config

import (
	"flag"
	"fmt"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	APIKey      string
	Concurrency int
	Verbose     bool
	InputFile   string
	InputURL    string
	MaxReviews  int
	OutputFile  string
	Language    string
}

// ParseFlags parses command-line flags and returns a Config struct
func ParseFlags() *Config {
	apiKey := flag.String("api", "", "API key for authentication (required)")
	concurrency := flag.Int("c", 5, "Number of concurrent workers")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	inputFile := flag.String("f", "", "Text file containing Goodreads URLs (one per line)")
	maxReviews := flag.Int("m", 100, "Maximum number of reviews to scrape per book")
	outputFile := flag.String("o", "", "Output CSV file (default: auto-generated with timestamp)")
	language := flag.String("l", "id", "Language code for reviews")

	flag.Parse()

	// Check for positional argument (single URL)
	var url string
	if len(flag.Args()) > 0 {
		url = flag.Args()[0]
	}

	cfg := &Config{
		APIKey:      *apiKey,
		Concurrency: *concurrency,
		Verbose:     *verbose,
		InputFile:   *inputFile,
		InputURL:    url,
		MaxReviews:  *maxReviews,
		OutputFile:  *outputFile,
		Language:    *language,
	}

	// Set default output file if not provided
	if cfg.OutputFile == "" {
		timestamp := time.Now().Format("20060102_150405")
		cfg.OutputFile = fmt.Sprintf("results/goodreads_reviews_%s.csv", timestamp)
	}

	return cfg
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required. Use -api flag to provide it")
	}
	return nil
}
