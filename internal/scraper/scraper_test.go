package scraper

import (
	"encoding/json"
	"goodreadscrape/internal/models"
	"testing"
)

func TestGraphQLResponseParsing(t *testing.T) {
	// Sample GraphQL response based on the actual format
	sampleResponse := `{
		"data": {
			"getReviews": {
				"totalCount": 1196,
				"edges": [
					{
						"node": {
							"__typename": "Review",
							"id": "kca://review:goodreads/amzn1.gr.review:goodreads.v1.yWCxqUxA5Gl5ZPYmYoRbkQ",
							"creator": {
								"id": 10113893,
								"imageUrlSquare": "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/users/1339307190i/10113893._UX200_CR0,35,200,200_.jpg",
								"isAuthor": true,
								"viewerRelationshipStatus": null,
								"followersCount": 7,
								"__typename": "User",
								"textReviewsCount": 54,
								"name": "Erma",
								"webUrl": "https://www.goodreads.com/user/show/10113893-erma",
								"contributor": {
									"id": "kca://author/amzn1.gr.author.v1.HzGNYiOaf6e58Kyt-z9s_w",
									"works": {
										"totalCount": 1,
										"__typename": "ContributorWorksConnection"
									},
									"__typename": "Contributor"
								}
							},
							"recommendFor": null,
							"updatedAt": 1761546360777,
							"createdAt": 1390843144000,
							"spoilerStatus": false,
							"lastRevisionAt": 1390843257000,
							"text": "Jika aku pernah bilang novel yang bagus adalah novel yang tidak bisa membuat tidur pembacanya, maka lain dengan novel ini.",
							"rating": 5,
							"shelving": {
								"shelf": {
									"name": "read",
									"displayName": "Read",
									"editable": false,
									"default": true,
									"actionType": null,
									"sortOrder": null,
									"webUrl": "https://www.goodreads.com/review/list/10113893?shelf=read",
									"__typename": "Shelf"
								},
								"taggings": [],
								"webUrl": "https://www.goodreads.com/review/show/836115237",
								"__typename": "Shelving"
							},
							"likeCount": 80,
							"viewerHasLiked": null,
							"commentCount": 16
						},
						"__typename": "BookReviewsEdge"
					}
				],
				"pageInfo": {
					"prevPageToken": "ODAwMywxMzkwODQzMTQ0MDAw",
					"nextPageToken": "NzEzNCwxNTE5OTIyNTQyMjMx",
					"__typename": "PageInfo"
				},
				"__typename": "BookReviewsConnection"
			}
		},
		"errors": [
			{
				"path": ["getReviews", "edges", 0, "node", "viewerHasLiked"],
				"data": null,
				"errorType": "Unauthorized",
				"errorInfo": null,
				"locations": [
					{
						"line": 63,
						"column": 11,
						"sourceName": null
					}
				],
				"message": "Not Authorized to access viewerHasLiked on type Review"
			}
		]
	}`

	// Test JSON parsing
	var graphqlResp GraphQLResponse
	err := json.Unmarshal([]byte(sampleResponse), &graphqlResp)
	if err != nil {
		t.Fatalf("Failed to parse GraphQL response: %v", err)
	}

	// Verify response structure
	if graphqlResp.Data.GetReviews.TotalCount != 1196 {
		t.Errorf("Expected totalCount 1196, got %d", graphqlResp.Data.GetReviews.TotalCount)
	}

	if len(graphqlResp.Data.GetReviews.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(graphqlResp.Data.GetReviews.Edges))
	}

	if len(graphqlResp.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(graphqlResp.Errors))
	}

	// Test review node parsing
	edge := graphqlResp.Data.GetReviews.Edges[0]
	node := edge.Node

	if node.ID != "kca://review:goodreads/amzn1.gr.review:goodreads.v1.yWCxqUxA5Gl5ZPYmYoRbkQ" {
		t.Errorf("Unexpected review ID: %s", node.ID)
	}

	if node.Creator.Name != "Erma" {
		t.Errorf("Expected creator name 'Erma', got '%s'", node.Creator.Name)
	}

	if node.Rating != 5 {
		t.Errorf("Expected rating 5, got %d", node.Rating)
	}

	if node.CreatedAt != 1390843144000 {
		t.Errorf("Expected createdAt 1390843144000, got %f", node.CreatedAt)
	}

	// Test error handling for unauthorized errors
	if graphqlResp.Errors[0].ErrorType != "Unauthorized" {
		t.Errorf("Expected error type 'Unauthorized', got '%s'", graphqlResp.Errors[0].ErrorType)
	}
}

func TestExtractReviewFromGraphQL(t *testing.T) {
	scraper := &goodreadsScraper{apiKey: "test"}

	// Sample review node
	node := ReviewNode{
		ID:        "test-review-123",
		Text:      "This is a <b>great</b> book with <br/>excellent content.",
		Rating:    4,
		CreatedAt: 1390843144000, // 2014-01-27/28 in milliseconds (depends on timezone)
		Creator: struct {
			ID                       int         `json:"id"`
			Name                     string      `json:"name"`
			ImageURLSquare           string      `json:"imageUrlSquare"`
			WebURL                   string      `json:"webUrl"`
			IsAuthor                 bool        `json:"isAuthor"`
			TextReviewsCount         int         `json:"textReviewsCount"`
			FollowersCount           int         `json:"followersCount"`
			ViewerRelationshipStatus interface{} `json:"viewerRelationshipStatus"`
			Contributor              interface{} `json:"contributor"`
			Typename                 string      `json:"__typename"`
		}{
			ID:   12345,
			Name: "Test Reviewer",
		},
	}

	bookMetadata := models.BookMetadata{
		Title:  "Test Book",
		Author: "Test Author",
		URL:    "https://www.goodreads.com/book/show/123",
	}

	review := scraper.extractReviewFromGraphQL(node, bookMetadata)

	// Test extracted data
	if review.ReviewID != "test-review-123" {
		t.Errorf("Expected ReviewID 'test-review-123', got '%s'", review.ReviewID)
	}

	if review.ReviewerName != "Test Reviewer" {
		t.Errorf("Expected ReviewerName 'Test Reviewer', got '%s'", review.ReviewerName)
	}

	if review.Rating != "4" {
		t.Errorf("Expected Rating '4', got '%s'", review.Rating)
	}

	if review.BookTitle != "Test Book" {
		t.Errorf("Expected BookTitle 'Test Book', got '%s'", review.BookTitle)
	}

	if review.BookURL != "https://www.goodreads.com/book/show/123" {
		t.Errorf("Expected BookURL 'https://www.goodreads.com/book/show/123', got '%s'", review.BookURL)
	}

	// Test HTML tag removal
	expectedText := "This is a great book with excellent content."
	if review.ReviewText != expectedText {
		t.Errorf("Expected ReviewText '%s', got '%s'", expectedText, review.ReviewText)
	}

	// Test date formatting (allow for timezone differences)
	expectedDates := []string{"2014-01-27", "2014-01-28"}
	dateFound := false
	for _, expectedDate := range expectedDates {
		if review.ReviewDate == expectedDate {
			dateFound = true
			break
		}
	}
	if !dateFound {
		t.Errorf("Expected ReviewDate to be one of %v, got '%s'", expectedDates, review.ReviewDate)
	}
}

func TestMinFunction(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{5, 3, 3},
		{1, 10, 1},
		{7, 7, 7},
		{0, 100, 0},
		{-5, -3, -5},
	}

	for _, test := range tests {
		result := min(test.a, test.b)
		if result != test.expected {
			t.Errorf("min(%d, %d) = %d; expected %d", test.a, test.b, result, test.expected)
		}
	}
}
