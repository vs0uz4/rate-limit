package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedisClientSuccess(t *testing.T) {
	client, err := NewRedisClient("localhost", "6379")
	assert.NoError(t, err, "Expected no error for successful Redis connection")
	assert.NotNil(t, client, "Expected client to be initialized")
}

func TestNewRedisClientFailure(t *testing.T) {
	client, err := NewRedisClient("localhost", "9999") // Porta incorreta
	assert.Error(t, err, "Expected error for failed Redis connection")
	assert.Nil(t, client, "Expected client to be nil on failure")
}
