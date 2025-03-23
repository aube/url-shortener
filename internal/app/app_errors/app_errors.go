package app_errors

import (
	"errors"
	"fmt"
)

type HTTPError struct {
	Code int
	Err  error
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("%d - %s", he.Code, he.Err)
}

func NewHTTPError(code int, message string) error {
	return &HTTPError{
		Code: code,
		Err:  errors.New(message),
	}
}
