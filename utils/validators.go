package utils

import (
	"fmt"
	"regexp"
)

var (
	idRegex     = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	emailRegex  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	urlRegex    = regexp.MustCompile(`^https?://[^\s]+$`)
)

// IsValidID checks if a string is a valid ID format.
func IsValidID(id string) bool {
	return id != "" && idRegex.MatchString(id)
}

// ValidateID returns an error if the ID is invalid.
func ValidateID(id, fieldName string) error {
	if !IsValidID(id) {
		return fmt.Errorf("%s must be a valid ID (alphanumeric, hyphens, underscores)", fieldName)
	}
	return nil
}

// IsValidEmail checks if a string is a valid email address.
func IsValidEmail(email string) bool {
	return email != "" && emailRegex.MatchString(email)
}

// ValidateEmail returns an error if the email is invalid.
func ValidateEmail(email, fieldName string) error {
	if !IsValidEmail(email) {
		return fmt.Errorf("%s must be a valid email address", fieldName)
	}
	return nil
}

// IsValidURL checks if a string is a valid URL.
func IsValidURL(urlStr string) bool {
	return urlRegex.MatchString(urlStr)
}

// ValidateNonNegative ensures the value is non-negative.
func ValidateNonNegative(value int, fieldName string) error {
	if value < 0 {
		return fmt.Errorf("%s must be non-negative", fieldName)
	}
	return nil
}
