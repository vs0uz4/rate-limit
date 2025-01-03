package errors

import "errors"

var (
	ErrMissingEnvFile = errors.New("failed to load .env file")
)
