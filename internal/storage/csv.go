package storage

import (
	"encoding/csv"
	"fmt"
	"goodreadscrape/internal/models"
	"os"
	"path/filepath"
)

// CSVStorage implements Storage interface for CSV output
type CSVStorage struct{}

// NewCSVStorage creates a new CSV storage instance
func NewCSVStorage() Storage {
	return &CSVStorage{}
}

// SaveReviews saves reviews to a CSV file
func (s *CSVStorage) SaveReviews(reviews []models.Review, outputPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Check if file exists to determine if we need to write header
	fileExists := false
	if _, err := os.Stat(outputPath); err == nil {
		fileExists = true
	}

	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header only if file is new
	if !fileExists {
		header := []string{
			"BookURL", "BookTitle", "ReviewID", "ReviewerName",
			"Rating", "ReviewText", "ReviewDate", "Language",
		}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	// Write reviews
	for _, review := range reviews {
		record := []string{
			review.BookURL,
			review.BookTitle,
			review.ReviewID,
			review.ReviewerName,
			review.Rating,
			review.ReviewText,
			review.ReviewDate,
			review.Language,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// SaveBookData saves book data to a CSV file
func (s *CSVStorage) SaveBookData(bookData models.BookData, outputPath string) error {
	// This is a placeholder implementation
	// You may want to implement specific logic for saving book data
	return fmt.Errorf("SaveBookData not yet implemented")
}
