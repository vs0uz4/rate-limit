package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) *RedisAdapter {
	return &RedisAdapter{client: client}
}

func (r *RedisAdapter) Incr(ctx context.Context, key string) (int64, error) {
	cmd := r.client.Incr(ctx, key)
	return cmd.Result()
}

func (r *RedisAdapter) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	cmd := r.client.SetNX(ctx, key, value, expiration)
	return cmd.Result()
}

func (r *RedisAdapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	cmd := r.client.TTL(ctx, key)
	return cmd.Result()
}

func (r *RedisAdapter) Expire(ctx context.Context, key string, expiration time.Duration) error {
	cmd := r.client.Expire(ctx, key, expiration)
	return cmd.Err()
}
