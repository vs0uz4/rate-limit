package mock

import (
	"context"
	"time"
)

type MockPersistenceProvider struct {
	IncrResponses map[string]int64
	IncrErrors    map[string]error
	SetNXErrors   map[string]error
	TTLResponses  map[string]time.Duration
	TTLErrors     map[string]error
	ExpireErrors  map[string]error
}

func NewMockPersistenceProvider() *MockPersistenceProvider {
	return &MockPersistenceProvider{
		IncrResponses: make(map[string]int64),
		IncrErrors:    make(map[string]error),
		SetNXErrors:   make(map[string]error),
		TTLResponses:  make(map[string]time.Duration),
		TTLErrors:     make(map[string]error),
		ExpireErrors:  make(map[string]error),
	}
}

func (m *MockPersistenceProvider) Incr(ctx context.Context, key string) (int64, error) {
	if err, exists := m.IncrErrors[key]; exists {
		return 0, err
	}
	m.IncrResponses[key]++
	return m.IncrResponses[key], nil
}

func (m *MockPersistenceProvider) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	if err, exists := m.SetNXErrors[key]; exists {
		return false, err
	}
	return true, nil
}

func (m *MockPersistenceProvider) TTL(ctx context.Context, key string) (time.Duration, error) {
	if err, exists := m.TTLErrors[key]; exists {
		return time.Duration(0), err
	}
	if ttl, exists := m.TTLResponses[key]; exists {
		return ttl, nil
	}
	return time.Duration(0), nil
}

func (m *MockPersistenceProvider) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err, exists := m.ExpireErrors[key]; exists {
		return err
	}
	return nil
}
