package storage

import (
	"encoding/csv"
	"goodreadscrape/internal/models"
	"os"
	"testing"
)

func TestNewCSVStorage(t *testing.T) {
	s := NewCSVStorage()
	if s == nil {
		t.Error("NewCSVStorage returned nil")
	}
	if _, ok := s.(*CSVStorage); !ok {
		t.Error("NewCSVStorage did not return *CSVStorage")
	}
}

func TestCSVStorage_SaveReviews(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "test_reviews_*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	tmpfile.Close()

	s := NewCSVStorage()
	reviews := []models.Review{
		{
			BookURL:      "http://example.com/book1",
			BookTitle:    "Test Book 1",
			ReviewID:     "123",
			ReviewerName: "John Doe",
			Rating:       "5",
			ReviewText:   "Great book!",
			ReviewDate:   "2023-01-01",
			Language:     "en",
		},
		{
			BookURL:      "http://example.com/book2",
			BookTitle:    "Test Book 2",
			ReviewID:     "456",
			ReviewerName: "Jane Doe",
			Rating:       "4",
			ReviewText:   "Good book.",
			ReviewDate:   "2023-01-02",
			Language:     "fr",
		},
	}

	// Test successful save
	err = s.SaveReviews(reviews, tmpfile.Name())
	if err != nil {
		t.Errorf("SaveReviews failed: %v", err)
	}

	// Verify file content
	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	// Check header
	expectedHeader := []string{
		"BookURL", "BookTitle", "ReviewID", "ReviewerName",
		"Rating", "ReviewText", "ReviewDate", "Language",
	}
	if len(records) < 1 {
		t.Fatal("File is empty")
	}
	for i, h := range expectedHeader {
		if records[0][i] != h {
			t.Errorf("Header mismatch at index %d: expected %s, got %s", i, h, records[0][i])
		}
	}

	// Check records
	if len(records) != 3 { // Header + 2 reviews
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	// Verify first review
	if records[1][0] != reviews[0].BookURL {
		t.Errorf("Expected BookURL %s, got %s", reviews[0].BookURL, records[1][0])
	}
	if records[1][1] != reviews[0].BookTitle {
		t.Errorf("Expected BookTitle %s, got %s", reviews[0].BookTitle, records[1][1])
	}

	// Test invalid path
	err = s.SaveReviews(reviews, "/invalid/path/to/file.csv")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestCSVStorage_SaveBookData(t *testing.T) {
	s := NewCSVStorage()
	err := s.SaveBookData(models.BookData{}, "dummy.csv")
	if err == nil {
		t.Error("Expected error from SaveBookData, got nil")
	}
	if err.Error() != "SaveBookData not yet implemented" {
		t.Errorf("Expected 'SaveBookData not yet implemented', got '%v'", err)
	}
}
