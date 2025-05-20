package apperrors

import (
	"errors"
	"fmt"
)

// HTTPError represents an HTTP error with status code and message.
type HTTPError struct {
	Code int   // HTTP status code
	Err  error // Error message
}

// Error implements the error interface for HTTPError.
func (he *HTTPError) Error() string {
	return fmt.Sprintf("%d - %s", he.Code, he.Err)
}

// NewHTTPError creates a new HTTPError with the given status code and message.
func NewHTTPError(code int, message string) error {
	return &HTTPError{
		Code: code,
		Err:  errors.New(message),
	}
}
