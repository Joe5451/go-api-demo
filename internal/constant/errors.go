package constant

import "errors"

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
)

const (
	ErrValidationCode = "VALIDATION_ERROR"
	ErrNotFoundCode   = "NOT_FOUND"
)
