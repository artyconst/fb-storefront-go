package review

import (
	"context"
	"fmt"
	"net/url"

	sf "github.com/artyconst/fb-storefront-go"
)

// ReviewService handles review-related operations.
type ReviewService struct {
	client *sf.StorefrontClient
}

// NewReviewService creates a new Review service instance.
func NewReviewService(client *sf.StorefrontClient) *ReviewService {
	return &ReviewService{client: client}
}

// CountByStore returns the total count of reviews for a specific store.
func (s *ReviewService) CountByStore(ctx context.Context, storeID string) (int, error) {
	if storeID == "" {
		return 0, fmt.Errorf("store ID cannot be empty")
	}

	path := "/reviews/count?store=" + storeID

	var result struct {
		Count int `json:"count"`
	}

	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return 0, fmt.Errorf("failed to get review count for store: %w", err)
	}

	return result.Count, nil
}

// CountByRating returns the total count of reviews with a specific rating.
func (s *ReviewService) CountByRating(ctx context.Context, rating int) (int, error) {
	if rating < 1 || rating > 5 {
		return 0, fmt.Errorf("rating must be between 1 and 5")
	}

	path := "/reviews/count?rating=" + fmt.Sprint(rating)

	var result struct {
		Count int `json:"count"`
	}

	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return 0, fmt.Errorf("failed to get review count for rating: %w", err)
	}

	return result.Count, nil
}

// List retrieves all reviews - GET /reviews
func (s *ReviewService) List(ctx context.Context, opts ...ListOption) ([]*Review, error) {
	options := &ListOptions{}
	for _, o := range opts {
		o(options)
	}

	path := "/reviews"
	if options.Limit > 0 || options.Offset > 0 {
		params := make(url.Values)
		if options.Limit > 0 {
			params.Set("limit", fmt.Sprintf("%d", options.Limit))
		}
		if options.Offset > 0 {
			params.Set("offset", fmt.Sprintf("%d", options.Offset))
		}
		path += "?" + params.Encode()
	}

	var reviews []*Review
	if err := s.client.GetJSON(ctx, path, &reviews); err != nil {
		return nil, fmt.Errorf("failed to list reviews: %w", err)
	}
	return reviews, nil
}

// Create creates a new review - POST /reviews
func (s *ReviewService) Create(ctx context.Context, content string, rating int64, opts ...CreateOption) (*Review, error) {
	if content == "" {
		return nil, fmt.Errorf("review content is required")
	}
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	options := &CreateReviewOptions{Content: content, Rating: int(rating)}
	for _, o := range opts {
		o(options)
	}

	var review Review
	if err := s.client.PostJSON(ctx, "/reviews", options, &review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}
	return &review, nil
}

// Get retrieves a single review by ID - GET /reviews/{id}
func (s *ReviewService) Get(ctx context.Context, id string) (*Review, error) {
	if id == "" {
		return nil, fmt.Errorf("review ID is required")
	}

	var review Review
	endpoint := "/reviews/" + id
	if err := s.client.GetJSON(ctx, endpoint, &review); err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}
	return &review, nil
}

// Delete deletes a review by ID - DELETE /reviews/{id}
func (s *ReviewService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("review ID is required")
	}

	endpoint := "/reviews/" + id
	if err := s.client.DeleteJSON(ctx, endpoint, nil); err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}

// GetMedia returns the media URLs associated with this review.
func (r *Review) GetMedia() []string {
	if r == nil {
		return nil
	}
	return r.MediaURLs
}
