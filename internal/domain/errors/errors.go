package errors

import "errors"

var (
	ErrGettingExecPath   = errors.New("configuration error: failed to get executable path")
	ErrEnvFileNotFound   = errors.New("configuration error: .env file not found or could not be read")
	ErrRedisHostRequired = errors.New("configuration error: REDIS_HOST is required")
)
