package rate_limiter

import (
	"context"
	"time"

	"github.com/vs0uz4/rate-limit/internal/contract"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

type Config struct {
	Limit         int
	BlockDuration time.Duration
}

type RedisRateLimiter struct {
	config Config
	redis  contract.RedisClient
}

func NewRedisRateLimiter(config Config, redis contract.RedisClient) *RedisRateLimiter {
	return &RedisRateLimiter{
		config: config,
		redis:  redis,
	}
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	count, err := r.redis.Incr(ctx, key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		_, err = r.redis.SetNX(ctx, key, 1, r.config.BlockDuration)
		if err != nil {
			return false, err
		}
	}

	if count > int64(r.config.Limit) {
		return false, nil
	}

	return true, nil
}
