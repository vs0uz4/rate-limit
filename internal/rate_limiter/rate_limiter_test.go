package rate_limiter

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vs0uz4/rate-limit/internal/mock"
)

func TestRateLimiterTokenConfigured(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
		TokenLimits: map[string]int{
			"token123": 10,
		},
	}, mockStore)

	for i := 0; i < 10; i++ {
		allowed, err := limiter.Allow(context.Background(), "test_key", "token123")
		assert.NoError(t, err)
		assert.True(t, allowed, "Request with configured token must be allowed")
	}

	allowed, err := limiter.Allow(context.Background(), "test_key", "token123")
	assert.NoError(t, err)
	assert.False(t, allowed, "Request with configured token must be blocked if it exceeds the limit")
}

func TestRateLimiterTokenNotConfigured(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
		TokenLimits:   map[string]int{},
	}, mockStore)

	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(context.Background(), "test_key", "unknown_token")
		assert.NoError(t, err)
		assert.True(t, allowed, "Request without token configuration must use global limit")
	}

	allowed, err := limiter.Allow(context.Background(), "test_key", "unknown_token")
	assert.NoError(t, err)
	assert.False(t, allowed, "Request without token configuration should be blocked when exceeding global limit")
}

func TestRateLimiterNoToken(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
	}, mockStore)

	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(context.Background(), "test_key", "")
		assert.NoError(t, err)
		assert.True(t, allowed, "Request without token must use global limit")
	}

	allowed, err := limiter.Allow(context.Background(), "test_key", "")
	assert.NoError(t, err)
	assert.False(t, allowed, "Request without token must be blocked if it exceeds global limit")
}

func TestAllowSuccess(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
		TokenLimits:   map[string]int{},
	}, mockStore)

	for i := 0; i < 3; i++ {
		allowed, err := limiter.Allow(context.Background(), "test_key", "")
		assert.NoError(t, err)
		assert.True(t, allowed, "Request must be allowed")
	}
}

func TestAllowLimitExceeded(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	limiter := NewRateLimiter(Config{
		Limit:         3,
		BlockDuration: time.Minute,
		TokenLimits: map[string]int{
			"token123": 10,
		},
	}, mockStore)

	for i := 0; i < 4; i++ {
		allowed, err := limiter.Allow(context.Background(), "test_key", "token456")
		if i < 3 {
			assert.NoError(t, err)
			assert.True(t, allowed, "Request must be allowed")
		} else {
			assert.NoError(t, err)
			assert.False(t, allowed, "Request must be blocked")
		}
	}
}

func TestRateLimiterExpire(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	rateLimiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
	}, mockStore)

	t.Run("Should set TTL when key exists without expiration", func(t *testing.T) {
		mockStore.TTLResponses = map[string]time.Duration{
			"ip:127.0.0.1:rate": 0,
		}

		mockStore.ExpireErrors = make(map[string]error)

		allowed, err := rateLimiter.Allow(context.Background(), "127.0.0.1", "")
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Should return error if Expire fails", func(t *testing.T) {
		mockStore.TTLResponses = map[string]time.Duration{
			"ip:127.0.0.1:rate": 0,
		}

		mockStore.ExpireErrors = map[string]error{
			"ip:127.0.0.1:rate": errors.New("expire failed"),
		}

		allowed, err := rateLimiter.Allow(context.Background(), "127.0.0.1", "")
		assert.Error(t, err)
		assert.False(t, allowed)
		assert.Equal(t, "expire failed", err.Error())
	})
}

func TestAllowErrorOnIncr(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	mockStore.IncrErrors["ip:127.0.0.1:rate"] = errors.New("error on incr")

	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
	}, mockStore)

	allowed, err := limiter.Allow(context.Background(), "127.0.0.1", "")
	assert.Error(t, err, "Expected increment error")
	assert.False(t, allowed, "Request must be blocked")
	assert.Equal(t, "error on incr", err.Error())
}

func TestAllowErrorOnTTLRateKey(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	mockStore.TTLErrors["ip:127.0.0.1:rate"] = errors.New("error on TTL for rateKey")

	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
	}, mockStore)

	allowed, err := limiter.Allow(context.Background(), "127.0.0.1", "")
	assert.Error(t, err, "Expected error when checking TTL for rateKey")
	assert.False(t, allowed, "Request must be blocked due to TTL error on rateKey")
	assert.Equal(t, "error on TTL for rateKey", err.Error(), "Unexpected error message")
}

func TestAllowErrorOnExpireRateKey(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	mockStore.TTLResponses["ip:127.0.0.1:rate"] = time.Duration(0)
	mockStore.ExpireErrors["ip:127.0.0.1:rate"] = errors.New("error on expire")

	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
	}, mockStore)

	allowed, err := limiter.Allow(context.Background(), "127.0.0.1", "")
	assert.Error(t, err, "Expected error when expiring rateKey")
	assert.False(t, allowed, "Request must be blocked due to expire error on rateKey")
	assert.Equal(t, "error on expire", err.Error(), "Unexpected error message")
}

func TestAllowErrorOnTTLBlockKey(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	mockStore.TTLErrors["ip:127.0.0.1:block"] = errors.New("error on TTL for blockKey")

	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
	}, mockStore)

	allowed, err := limiter.Allow(context.Background(), "127.0.0.1", "")
	assert.Error(t, err, "Expected error when checking TTL for blockKey")
	assert.False(t, allowed, "Request must be blocked due to TTL error on blockKey")
	assert.Equal(t, "error on TTL for blockKey", err.Error(), "Unexpected error message")
}

func TestAllowErrorOnSetNX(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	mockStore.IncrResponses = map[string]int64{
		"ip:127.0.0.1:rate": 6,
	}
	mockStore.SetNXErrors = map[string]error{
		"ip:127.0.0.1:block": errors.New("error on setnx"),
	}

	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
	}, mockStore)

	fmt.Println("Running TestAllowErrorOnSetNX...")
	allowed, err := limiter.Allow(context.Background(), "127.0.0.1", "")
	fmt.Printf("Result - Allowed: %v, Error: %v\n", allowed, err)

	assert.Error(t, err, "Expected error when setting TTL")
	assert.False(t, allowed, "Request must be blocked")
	assert.Equal(t, "error on setnx", err.Error(), "Expected error when setting TTL")
}

func TestAllowBlockKeyTTLNotExpired(t *testing.T) {
	mockStore := mock.NewMockPersistenceProvider()
	mockStore.TTLResponses["ip:127.0.0.1:block"] = time.Minute
	mockStore.TTLResponses["ip:127.0.0.1:rate"] = time.Second * 1

	limiter := NewRateLimiter(Config{
		Limit:         5,
		BlockDuration: time.Minute,
	}, mockStore)

	allowed, err := limiter.Allow(context.Background(), "127.0.0.1", "")
	assert.NoError(t, err, "No error expected")
	assert.False(t, allowed, "Request must be blocked due to TTL not expired")
}
