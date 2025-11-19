package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractBookMetadata parses HTML document to extract book metadata
func ExtractBookMetadata(doc *goquery.Document, bookURL string) (string, string, float64) {
	// Extract book title
	title := "Unknown Title"
	titleElement := doc.Find("a[data-testid='title']").First()
	if titleElement.Length() > 0 {
		title = strings.TrimSpace(titleElement.Text())
	}

	// Extract author name
	author := "Unknown Author"
	authorElement := doc.Find("span.ContributorLink__name[data-testid='name']").First()
	if authorElement.Length() > 0 {
		author = strings.TrimSpace(authorElement.Text())
	}

	// Extract average rating
	avgRating := 0.0
	ratingElement := doc.Find("div.RatingStatistics__column").First()
	if ratingElement.Length() > 0 {
		// Try to get rating from aria-label first
		ariaLabel, exists := ratingElement.Attr("aria-label")
		if exists {
			re := regexp.MustCompile(`(\d+\.?\d*)`)
			matches := re.FindStringSubmatch(ariaLabel)
			if len(matches) > 1 {
				if rating, err := strconv.ParseFloat(matches[1], 64); err == nil {
					avgRating = rating
				}
			}
		}

		// Fallback to text content if aria-label didn't work
		if avgRating == 0.0 {
			ratingText := strings.TrimSpace(ratingElement.Text())
			re := regexp.MustCompile(`(\d+\.?\d*)`)
			matches := re.FindStringSubmatch(ratingText)
			if len(matches) > 1 {
				if rating, err := strconv.ParseFloat(matches[1], 64); err == nil {
					avgRating = rating
				}
			}
		}
	}

	return title, author, avgRating
}

