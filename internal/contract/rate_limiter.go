package contract

import "context"

type RateLimiter interface {
	Allow(ctx context.Context, key string, token string) (bool, error)
}
