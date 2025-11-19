package validator

import (
	"net/url"
	"strings"
)

// ValidateGoodreadsURL checks if the provided URL is a valid Goodreads book URL
func ValidateGoodreadsURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Check if the hostname matches Goodreads
	if parsed.Host != "www.goodreads.com" {
		return false
	}

	// Check if the path contains "/book/show/"
	return strings.Contains(parsed.Path, "/book/show/")
}

