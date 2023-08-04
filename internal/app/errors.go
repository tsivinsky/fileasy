package app

import "fmt"

type ApiError struct {
	Status  int
	Message string
	Err     *error
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

func NewApiError(status int, message string, err *error) *ApiError {
	return &ApiError{
		Status:  status,
		Message: message,
		Err:     err,
	}
}
