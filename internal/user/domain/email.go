package domain

import (
	"fmt"
	"regexp"
)

// Email represents a validated email address value object
type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object with validation
func NewEmail(value string) (Email, error) {
	if value == "" {
		return "", fmt.Errorf("email cannot be empty")
	}
	if !emailRegex.MatchString(value) {
		return "", fmt.Errorf("invalid email format: %s", value)
	}
	return Email(value), nil
}

// String returns the string representation
func (e Email) String() string {
	return string(e)
}

// IsValid checks if the email is valid
func (e Email) IsValid() bool {
	return emailRegex.MatchString(string(e))
}
