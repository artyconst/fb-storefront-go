package review

import "fmt"

// Review represents a store/product review.
type Review struct {
	ID        string   `json:"uid"`
	Rating    int      `json:"rating"`
	Content   string   `json:"content,omitempty"`
	MediaURLs []string `json:"media_urls,omitempty"`
	StoreID   string   `json:"store_uid,omitempty"`
}

// CountByRatingError represents an error related to rating count operations.
type CountByRatingError struct {
	Rating  int
	Message string
}

func (e *CountByRatingError) Error() string {
	return fmt.Sprintf("rating count failed for rating %d: %s", e.Rating, e.Message)
}

// NewCountByRatingError creates a new CountByRatingError.
func NewCountByRatingError(rating int, message string) *CountByRatingError {
	return &CountByRatingError{
		Rating:  rating,
		Message: message,
	}
}

// ListOptions holds query parameters for listing reviews.
type ListOptions struct {
	Limit  int64 `json:"limit,omitempty"`
	Offset int64 `json:"offset,omitempty"`
}

// ListOption is a functional option for configuring review list queries.
type ListOption func(*ListOptions)

// WithReviewLimit sets the limit for listing reviews.
func WithReviewLimit(limit int64) ListOption {
	return func(o *ListOptions) { o.Limit = limit }
}

// WithReviewOffset sets the offset for listing reviews.
func WithReviewOffset(offset int64) ListOption {
	return func(o *ListOptions) { o.Offset = offset }
}

// CreateReviewOptions holds parameters for creating a review.
type CreateReviewOptions struct {
	Content string `json:"content"`
	Rating  int    `json:"rating"`
}

// CreateOption is a functional option for configuring review creation.
type CreateOption func(*CreateReviewOptions)

// WithContent sets the content for a new review.
func WithContent(content string) CreateOption {
	return func(o *CreateReviewOptions) { o.Content = content }
}

// WithRating sets the rating for a new review.
func WithRating(rating int) CreateOption {
	return func(o *CreateReviewOptions) { o.Rating = rating }
}
