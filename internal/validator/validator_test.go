package validator

import "testing"

func TestValidateGoodreadsURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Valid Goodreads URL",
			url:      "https://www.goodreads.com/book/show/12345.Some_Book",
			expected: true,
		},
		{
			name:     "Invalid URL format",
			url:      "://invalid-url",
			expected: false,
		},
		{
			name:     "Incorrect Hostname",
			url:      "https://www.badreads.com/book/show/12345",
			expected: false,
		},
		{
			name:     "Missing /book/show/ in path",
			url:      "https://www.goodreads.com/author/show/12345",
			expected: false,
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: false,
		},
		{
			name:     "Valid URL with query params",
			url:      "https://www.goodreads.com/book/show/54321?ref=nav_som",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateGoodreadsURL(tt.url)
			if result != tt.expected {
				t.Errorf("ValidateGoodreadsURL(%q) = %v; want %v", tt.url, result, tt.expected)
			}
		})
	}
}
