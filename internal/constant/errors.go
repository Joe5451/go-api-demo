package constant

import "errors"

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
)

type ErrorCode string

const (
	ErrValidationCode      ErrorCode = "VALIDATION_ERROR"
	ErrNotFoundCode        ErrorCode = "NOT_FOUND"
	ErrInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
)
