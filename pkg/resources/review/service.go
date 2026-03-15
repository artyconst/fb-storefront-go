package review

import (
	"context"
	"fmt"

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

// GetMedia returns the media URLs associated with this review.
func (r *Review) GetMedia() []string {
	if r == nil {
		return nil
	}
	return r.MediaURLs
}
