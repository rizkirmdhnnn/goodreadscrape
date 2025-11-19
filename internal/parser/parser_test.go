package parser

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractBookMetadata(t *testing.T) {
	tests := []struct {
		name           string
		html           string
		bookURL        string
		expectedTitle  string
		expectedAuthor string
		expectedRating float64
	}{
		{
			name: "Valid HTML with all metadata",
			html: `
				<html>
					<body>
						<a data-testid="title">Test Book Title</a>
						<span class="ContributorLink__name" data-testid="name">Test Author</span>
						<div class="RatingStatistics__column" aria-label="4.5 out of 5 stars">4.5</div>
					</body>
				</html>
			`,
			bookURL:        "http://example.com/book",
			expectedTitle:  "Test Book Title",
			expectedAuthor: "Test Author",
			expectedRating: 4.5,
		},
		{
			name: "Missing title",
			html: `
				<html>
					<body>
						<span class="ContributorLink__name" data-testid="name">Test Author</span>
						<div class="RatingStatistics__column" aria-label="4.5 out of 5 stars">4.5</div>
					</body>
				</html>
			`,
			bookURL:        "http://example.com/book",
			expectedTitle:  "Unknown Title",
			expectedAuthor: "Test Author",
			expectedRating: 4.5,
		},
		{
			name: "Missing author",
			html: `
				<html>
					<body>
						<a data-testid="title">Test Book Title</a>
						<div class="RatingStatistics__column" aria-label="4.5 out of 5 stars">4.5</div>
					</body>
				</html>
			`,
			bookURL:        "http://example.com/book",
			expectedTitle:  "Test Book Title",
			expectedAuthor: "Unknown Author",
			expectedRating: 4.5,
		},
		{
			name: "Missing rating",
			html: `
				<html>
					<body>
						<a data-testid="title">Test Book Title</a>
						<span class="ContributorLink__name" data-testid="name">Test Author</span>
					</body>
				</html>
			`,
			bookURL:        "http://example.com/book",
			expectedTitle:  "Test Book Title",
			expectedAuthor: "Test Author",
			expectedRating: 0.0,
		},
		{
			name: "Rating in aria-label",
			html: `
				<html>
					<body>
						<a data-testid="title">Test Book Title</a>
						<span class="ContributorLink__name" data-testid="name">Test Author</span>
						<div class="RatingStatistics__column" aria-label="Rating 3.8"></div>
					</body>
				</html>
			`,
			bookURL:        "http://example.com/book",
			expectedTitle:  "Test Book Title",
			expectedAuthor: "Test Author",
			expectedRating: 3.8,
		},
		{
			name: "Rating in text content",
			html: `
				<html>
					<body>
						<a data-testid="title">Test Book Title</a>
						<span class="ContributorLink__name" data-testid="name">Test Author</span>
						<div class="RatingStatistics__column">4.2</div>
					</body>
				</html>
			`,
			bookURL:        "http://example.com/book",
			expectedTitle:  "Test Book Title",
			expectedAuthor: "Test Author",
			expectedRating: 4.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			title, author, rating := ExtractBookMetadata(doc, tt.bookURL)

			if title != tt.expectedTitle {
				t.Errorf("Expected title %q, got %q", tt.expectedTitle, title)
			}
			if author != tt.expectedAuthor {
				t.Errorf("Expected author %q, got %q", tt.expectedAuthor, author)
			}
			if rating != tt.expectedRating {
				t.Errorf("Expected rating %f, got %f", tt.expectedRating, rating)
			}
		})
	}
}
