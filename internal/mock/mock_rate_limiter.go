package mock

import (
	"context"
)

type MockRateLimiter struct {
	Allowed bool
	Err     error
}

func (m *MockRateLimiter) Allow(ctx context.Context, key string, token string) (bool, error) {
	return m.Allowed, m.Err
}
