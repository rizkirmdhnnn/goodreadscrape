package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goodreadscrape/internal/models"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Define GoodreadsScraper interface and its implementation
type goodreadsScraper struct {
	apiKey  string
	verbose bool
}

type GoodreadsScraper interface {
	ScrapeBookData(bookURL string, maxReviews int, filters models.Filters) (models.BookData, error)
	ExtractBookMetadata(bookURL string) (models.BookMetadata, error)
	ExtractWorkID(bookURL string) (string, error)
	FetchReviewsGraphQL(workID string, maxReviews int, languageCode string, bookMetadata models.BookMetadata) ([]models.Review, error)
}

func NewGoodreadsScraper(apiKey string, verbose bool) GoodreadsScraper {
	return &goodreadsScraper{
		apiKey:  apiKey,
		verbose: verbose,
	}
}

// Implement the methods of GoodreadsScraper interface here
func (s *goodreadsScraper) ScrapeBookData(bookURL string, maxReviews int, filters models.Filters) (models.BookData, error) {
	// Extract Metadata
	metadata, err := s.ExtractBookMetadata(bookURL)
	if err != nil {
		return models.BookData{}, err
	}

	// Extract Work ID
	workID, err := s.ExtractWorkID(bookURL)
	if err != nil {
		return models.BookData{}, err
	}

	// Fetch Reviews using GraphQL API
	reviews, err := s.FetchReviewsGraphQL(workID, maxReviews, filters.Language, metadata)
	if err != nil {
		fmt.Printf("⚠️ Warning: Failed to fetch reviews via GraphQL: %v\n", err)
		return models.BookData{
			Metadata: metadata,
		}, nil
	}

	if s.verbose {
		fmt.Printf("✅ Successfully fetched %d reviews\n", len(reviews))
	}

	return models.BookData{
		Metadata: metadata,
		Reviews:  reviews,
	}, nil
}

// Implement the methods of GoodreadsScraper interface here
func (s *goodreadsScraper) ExtractBookMetadata(bookURL string) (models.BookMetadata, error) {
	// Convert bookURL to review URL
	ReviewURL := bookURL + "/reviews"

	// Request to url
	res, err := http.Get(ReviewURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Parse response
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Extract book title from the link
	title := "Unknown Title"
	titleElement := doc.Find("a[data-testid='title']").First()
	if titleElement.Length() > 0 {
		title = strings.TrimSpace(titleElement.Text())
	}

	// Extract author name from the link
	author := "Unknown Author"
	authorElement := doc.Find("span.ContributorLink__name[data-testid='name']").First()
	if authorElement.Length() > 0 {
		author = strings.TrimSpace(authorElement.Text())
	}

	// Extract average rating from the link
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

	return models.BookMetadata{
		Title:         title,
		Author:        author,
		AverageRating: avgRating,
		URL:           bookURL,
	}, nil
}

func (s *goodreadsScraper) ExtractWorkID(bookURL string) (string, error) {
	// Convert to reviews URL
	reviewsURL := strings.TrimSuffix(bookURL, "/") + "/reviews"

	// Create HTTP request with headers
	client := &http.Client{}
	req, err := http.NewRequest("GET", reviewsURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set User-Agent header to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Look for work ID in JavaScript using regex
	workIDRegex := regexp.MustCompile(`"work":\s*{\s*"__ref":\s*"Work:(kca://work/[^"]+)"`)
	matches := workIDRegex.FindStringSubmatch(string(body))

	if len(matches) > 1 {
		workID := matches[1]
		if s.verbose {
			fmt.Printf("✅ Found work ID: %s\n", workID)
		}
		return workID, nil
	}

	fmt.Println("❌ Work ID not found in JavaScript")
	return "", fmt.Errorf("work ID not found in page content")
}

// GraphQL types for API responses
type GraphQLResponse struct {
	Data struct {
		GetReviews struct {
			TotalCount int `json:"totalCount"`
			Edges      []struct {
				Node     ReviewNode `json:"node"`
				Typename string     `json:"__typename"`
			} `json:"edges"`
			PageInfo struct {
				PrevPageToken string `json:"prevPageToken"`
				NextPageToken string `json:"nextPageToken"`
				Typename      string `json:"__typename"`
			} `json:"pageInfo"`
			Typename string `json:"__typename"`
		} `json:"getReviews"`
	} `json:"data"`
	Errors []struct {
		Path      []interface{} `json:"path"`
		Data      interface{}   `json:"data"`
		ErrorType string        `json:"errorType"`
		ErrorInfo interface{}   `json:"errorInfo"`
		Locations []struct {
			Line       int    `json:"line"`
			Column     int    `json:"column"`
			SourceName string `json:"sourceName"`
		} `json:"locations"`
		Message string `json:"message"`
	} `json:"errors"`
}

type ReviewNode struct {
	Typename       string  `json:"__typename"`
	ID             string  `json:"id"`
	Text           string  `json:"text"`
	Rating         int     `json:"rating"`
	CreatedAt      float64 `json:"createdAt"`      // Changed to float64 for epoch timestamp
	UpdatedAt      float64 `json:"updatedAt"`      // Changed to float64 for epoch timestamp
	LastRevisionAt float64 `json:"lastRevisionAt"` // Added missing field
	SpoilerStatus  bool    `json:"spoilerStatus"`  // Added missing field
	Creator        struct {
		ID                       int         `json:"id"` // Changed to int
		Name                     string      `json:"name"`
		ImageURLSquare           string      `json:"imageUrlSquare"`
		WebURL                   string      `json:"webUrl"`
		IsAuthor                 bool        `json:"isAuthor"`
		TextReviewsCount         int         `json:"textReviewsCount"`
		FollowersCount           int         `json:"followersCount"`
		ViewerRelationshipStatus interface{} `json:"viewerRelationshipStatus"` // Can be null
		Contributor              interface{} `json:"contributor"`              // Can be null or object
		Typename                 string      `json:"__typename"`
	} `json:"creator"`
	RecommendFor   interface{} `json:"recommendFor"` // Can be null
	LikeCount      int         `json:"likeCount"`
	CommentCount   int         `json:"commentCount"`
	ViewerHasLiked interface{} `json:"viewerHasLiked"` // Can be null due to auth errors
	Shelving       struct {
		Shelf struct {
			Name        string      `json:"name"`
			DisplayName string      `json:"displayName"`
			Editable    bool        `json:"editable"`
			Default     bool        `json:"default"`
			ActionType  interface{} `json:"actionType"`
			SortOrder   interface{} `json:"sortOrder"`
			WebURL      string      `json:"webUrl"`
			Typename    string      `json:"__typename"`
		} `json:"shelf"`
		Taggings []interface{} `json:"taggings"`
		WebURL   string        `json:"webUrl"`
		Typename string        `json:"__typename"`
	} `json:"shelving"`
}

type GraphQLRequest struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Query         string                 `json:"query"`
}

// FetchReviewsGraphQL fetches reviews using GraphQL API
func (s *goodreadsScraper) FetchReviewsGraphQL(workID string, maxReviews int, languageCode string, bookMetadata models.BookMetadata) ([]models.Review, error) {
	var reviews []models.Review
	var afterToken string
	limit := 100 // API limit per request

	graphqlURL := "https://kxbwmqov6jgg3daaamb744ycu4.appsync-api.us-east-1.amazonaws.com/graphql"

	// Get API key from environment variable
	apiKey := os.Getenv("GOODREADS_API_KEY")
	if apiKey == "" {
		apiKey = s.apiKey
	}
	if apiKey == "" || apiKey == "xxxxxx" {
		return nil, fmt.Errorf("API key not set! Please set the GOODREADS_API_KEY environment variable or provide it via constructor")
	}

	headers := map[string]string{
		"accept":             "*/*",
		"accept-language":    "en-US,en;q=0.9",
		"content-type":       "application/json",
		"x-api-key":          apiKey,
		"sec-ch-ua":          `"Google Chrome";v="141", "Not?A_Brand";v="8", "Chromium";v="141"`,
		"sec-ch-ua-mobile":   "?1",
		"sec-ch-ua-platform": `"Android"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "cross-site",
	}

	query := `
        query getReviews($filters: BookReviewsFilterInput!, $pagination: PaginationInput) {
          getReviews(filters: $filters, pagination: $pagination) {
            ...BookReviewsFragment
            __typename
          }
        }

        fragment BookReviewsFragment on BookReviewsConnection {
          totalCount
          edges {
            node {
              ...ReviewCardFragment
              __typename
            }
            __typename
          }
          pageInfo {
            prevPageToken
            nextPageToken
            __typename
          }
          __typename
        }

        fragment ReviewCardFragment on Review {
          __typename
          id
          creator {
            ...ReviewerProfileFragment
            __typename
          }
          recommendFor
          updatedAt
          createdAt
          spoilerStatus
          lastRevisionAt
          text
          rating
          shelving {
            shelf {
              name
              displayName
              editable
              default
              actionType
              sortOrder
              webUrl
              __typename
            }
            taggings {
              tag {
                name
                webUrl
                __typename
              }
              __typename
            }
            webUrl
            __typename
          }
          likeCount
          viewerHasLiked
          commentCount
        }

        fragment ReviewerProfileFragment on User {
          id: legacyId
          imageUrlSquare
          isAuthor
          ...SocialUserFragment
          textReviewsCount
          viewerRelationshipStatus {
            isBlockedByViewer
            __typename
          }
          name
          webUrl
          contributor {
            id
            works {
              totalCount
              __typename
            }
            __typename
          }
          __typename
        }

        fragment SocialUserFragment on User {
          viewerRelationshipStatus {
            isFollowing
            isFriend
            __typename
          }
          followersCount
          __typename
        }
        `

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if s.verbose {
		fmt.Printf("🚀 Starting GraphQL review fetch for work ID: %s\n", workID)
		fmt.Printf("📋 Target: %d reviews, Language: %s\n", maxReviews, languageCode)
	}

	for len(reviews) < maxReviews {
		filters := map[string]interface{}{
			"resourceType": "WORK",
			"resourceId":   workID,
		}

		// Add language filter if specified
		if languageCode != "" {
			filters["languageCode"] = languageCode
			if s.verbose {
				fmt.Printf("🌍 Filtering reviews by language: %s\n", languageCode)
			}
		}

		variables := map[string]interface{}{
			"filters": filters,
			"pagination": map[string]interface{}{
				"limit": min(limit, maxReviews-len(reviews)),
			},
		}

		if afterToken != "" {
			variables["pagination"].(map[string]interface{})["after"] = afterToken
		}

		payload := GraphQLRequest{
			OperationName: "getReviews",
			Variables:     variables,
			Query:         query,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return reviews, fmt.Errorf("error marshaling GraphQL payload: %v", err)
		}

		req, err := http.NewRequest("POST", graphqlURL, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return reviews, fmt.Errorf("error creating request: %v", err)
		}

		// Set headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Error fetching reviews from GraphQL: %v\n", err)
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("❌ HTTP error: %d %s\n", resp.StatusCode, resp.Status)
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Response body: %s\n", string(body))
			break
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("❌ Error reading response body: %v\n", err)
			break
		}

		// log the raw response for debugging
		if s.verbose {
			fmt.Printf("GraphQL Response: %s\n", string(body))
		}

		var graphqlResp GraphQLResponse
		if err := json.Unmarshal(body, &graphqlResp); err != nil {
			fmt.Printf("❌ Error parsing GraphQL response: %v\n", err)
			break
		}

		if len(graphqlResp.Errors) > 0 {
			// Check if errors are just authorization errors for non-critical fields
			criticalError := false
			authErrorCount := 0
			for _, err := range graphqlResp.Errors {
				if err.ErrorType != "Unauthorized" {
					criticalError = true
					break
				} else {
					authErrorCount++
				}
			}
			if criticalError {
				fmt.Printf("❌ Critical GraphQL errors: %v\n", graphqlResp.Errors)
				break
			} else {
				if s.verbose {
					fmt.Printf("⚠️ %d non-critical authorization errors (expected for public API)\n", authErrorCount)
				}
			}
		}

		if len(graphqlResp.Data.GetReviews.Edges) == 0 {
			fmt.Printf("📊 No more reviews found. Total available: %d, Fetched: %d\n",
				graphqlResp.Data.GetReviews.TotalCount, len(reviews))
			break
		}

		if s.verbose {
			fmt.Printf("📦 Processing batch of %d reviews...\n", len(graphqlResp.Data.GetReviews.Edges))
		}

		batchProcessed := 0
		for _, edge := range graphqlResp.Data.GetReviews.Edges {
			if len(reviews) >= maxReviews {
				break
			}

			reviewData := s.extractReviewFromGraphQL(edge.Node, bookMetadata)
			if reviewData.ReviewID != "" {
				reviews = append(reviews, reviewData)
				batchProcessed++
			} else {
				fmt.Printf("⚠️ Skipped review due to missing data (ID: %s)\n", edge.Node.ID)
			}
		}
		if s.verbose {
			fmt.Printf("✅ Processed %d reviews from this batch\n", batchProcessed)
		}

		// Check for next page
		afterToken = graphqlResp.Data.GetReviews.PageInfo.NextPageToken
		if afterToken == "" {
			break
		}

		if s.verbose {
			fmt.Printf("📊 Progress: %d/%d reviews fetched (%.1f%% complete)\n",
				len(reviews), maxReviews, float64(len(reviews))/float64(maxReviews)*100)
		}

		if afterToken != "" {
			if s.verbose {
				fmt.Printf("🔄 Continuing to next page...\n")
			}
			time.Sleep(1 * time.Second) // Rate limiting
		}
	}

	finalCount := min(len(reviews), maxReviews)
	if s.verbose {
		fmt.Printf("🎉 GraphQL fetch completed! Retrieved %d reviews total\n", finalCount)
	}
	return reviews[:finalCount], nil
}

// extractReviewFromGraphQL converts GraphQL review node to models.Review
func (s *goodreadsScraper) extractReviewFromGraphQL(node ReviewNode, bookMetadata models.BookMetadata) models.Review {
	ratingStr := ""
	if node.Rating > 0 {
		ratingStr = strconv.Itoa(node.Rating)
	}

	// Format date - convert from epoch timestamp
	reviewDate := ""
	if node.CreatedAt > 0 {
		// Convert from milliseconds epoch to time
		createdTime := time.Unix(int64(node.CreatedAt/1000), 0)
		reviewDate = createdTime.Format("2006-01-02")
	}

	// Clean up review text (remove HTML tags)
	reviewText := node.Text
	if reviewText != "" {
		// Simple HTML tag removal
		re := regexp.MustCompile(`<[^>]*>`)
		reviewText = re.ReplaceAllString(reviewText, "")
	}

	return models.Review{
		BookURL:      bookMetadata.URL,
		BookTitle:    bookMetadata.Title,
		ReviewID:     node.ID,
		ReviewerName: node.Creator.Name,
		Rating:       ratingStr,
		ReviewText:   reviewText,
		ReviewDate:   reviewDate,
		Language:     "", // Language detection could be added here if needed
	}
}

// Helper function for min (Go doesn't have built-in min for int)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
