package domain

import "errors"

// Domain-specific errors
var (
	ErrUserNotFound = errors.New("user not found")
)

// IsNotFoundError checks if an error is a not-found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}
