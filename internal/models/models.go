package models

// BookMetadata represents metadata about a book
type BookMetadata struct {
	Title         string
	Author        string
	AverageRating float64
	URL           string
}

// Filters contains filtering options for scraping
type Filters struct {
	Language string
}

// BookData contains complete book information including metadata and reviews
type BookData struct {
	Metadata BookMetadata
	Reviews  []Review
}

// Review represents a single book review
type Review struct {
	BookURL      string
	BookTitle    string
	ReviewID     string
	ReviewerName string
	Rating       string
	ReviewText   string
	ReviewDate   string
	Language     string
}
