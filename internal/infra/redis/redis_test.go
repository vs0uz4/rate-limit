package redis

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func isRedisAvailable() bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", "[::1]:6379", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func TestNewRedisClientSuccess(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Skipping test as Redis is not available")
	}

	client, err := NewRedisClient("localhost", "6379")
	assert.NoError(t, err, "Expected no error for successful Redis connection")
	assert.NotNil(t, client, "Expected client to be initialized")
}

func TestNewRedisClientFailure(t *testing.T) {
	client, err := NewRedisClient("localhost", "9999") // Porta incorreta
	assert.Error(t, err, "Expected error for failed Redis connection")
	assert.Nil(t, client, "Expected client to be nil on failure")
}
