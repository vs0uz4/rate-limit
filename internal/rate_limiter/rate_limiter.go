package rate_limiter

import (
	"context"
	"log"
	"time"

	"github.com/vs0uz4/rate-limit/internal/contract"
)

type Config struct {
	Limit         int
	BlockDuration time.Duration
	TokenLimits   map[string]int
}

type RateLimiter struct {
	config Config
	store  contract.PersistenceProvider
}

func NewRateLimiter(config Config, store contract.PersistenceProvider) *RateLimiter {
	return &RateLimiter{
		config: config,
		store:  store,
	}
}

func (r *RateLimiter) Allow(ctx context.Context, key string, token string) (bool, error) {
	var limit int

	if token != "" {
		if tokenLimit, exists := r.config.TokenLimits[token]; exists {
			limit = tokenLimit
			key = "token:" + token
		} else {
			limit = r.config.Limit
		}
	} else {
		limit = r.config.Limit
		key = "ip:" + key
	}

	rateKey := key + ":rate"

	count, err := r.store.Incr(ctx, rateKey)
	if err != nil {
		return false, err
	}

	ttl, err := r.store.TTL(ctx, rateKey)
	if err != nil {
		return false, err
	}

	log.Printf("Checking TTL for rateKey: %v", ttl)
	if ttl <= 0 {
		err = r.store.Expire(ctx, rateKey, time.Second)
		if err != nil {
			return false, err
		}
	}

	if count > int64(limit) {
		blockKey := key + ":block"
		_, err := r.store.SetNX(ctx, blockKey, 1, r.config.BlockDuration)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	blockKey := key + ":block"
	ttl, err = r.store.TTL(ctx, blockKey)
	if err != nil {
		return false, err
	}

	if ttl > 0 {
		return false, nil
	}

	return true, nil
}
