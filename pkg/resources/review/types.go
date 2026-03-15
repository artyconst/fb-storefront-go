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
