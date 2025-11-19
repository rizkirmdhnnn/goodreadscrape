package storage

import "goodreadscrape/internal/models"

// Storage interface defines methods for storing scraped data
type Storage interface {
	SaveReviews(reviews []models.Review, outputPath string) error
	SaveBookData(bookData models.BookData, outputPath string) error
}
