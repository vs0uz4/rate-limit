package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vs0uz4/rate-limit/internal/domain/errors"
)

func TestLoadConfig(t *testing.T) {
	envContent := `
REDIS_HOST=testhost
REDIS_PORT=1234
LIMITER_IP_LIMIT=10
LIMITER_TOKEN_LIMIT=20
BLOCK_DURATION=60
`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err, "Failed to create temporary .env file")
	defer os.Remove(".env")

	config, err := LoadConfig()

	assert.NoError(t, err, "Unexpected error when loading configuration")
	assert.NotNil(t, config, "Configuration should not be null")
	assert.Equal(t, "testhost", config.RedisHost)
	assert.Equal(t, "1234", config.RedisPort)
	assert.Equal(t, 10, config.LimiterIPLimit)
	assert.Equal(t, 20, config.LimiterTokenLimit)
	assert.Equal(t, 60, config.BlockDuration)
}

func TestLoadConfigMissingEnvFile(t *testing.T) {
	os.Remove(".env")

	config, err := LoadConfig()

	assert.Nil(t, config, "Configuration should be null if .env is missing")
	assert.Error(t, err, "Expected an error when the .env file was missing")
	assert.Equal(t, errors.ErrEnvFileNotFound, err, "Unexpected error message")
}

func TestLoadConfigMissingRedisHost(t *testing.T) {
	envContent := `
REDIS_PORT=1234
LIMITER_IP_LIMIT=10
LIMITER_TOKEN_LIMIT=20
BLOCK_DURATION=60
`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err, "Failed to create temporary .env file")
	defer os.Remove(".env")

	config, err := LoadConfig()

	assert.Nil(t, config, "Configuration should be null if REDIS_HOST is missing")
	assert.Error(t, err, "Expected an error when REDIS_HOST is missing")
	assert.Equal(t, errors.ErrRedisHostRequired, err, "Unexpected error message")
}

func TestLoadConfigWithDefaultValues(t *testing.T) {
	envContent := `
REDIS_HOST=testhost
`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err, "Failed to create temporary .env file")
	defer os.Remove(".env")

	config, err := LoadConfig()

	assert.NoError(t, err, "Unexpected error when loading configuration")
	assert.NotNil(t, config, "Configuration should not be null")
	assert.Equal(t, "testhost", config.RedisHost)
	assert.Equal(t, "6379", config.RedisPort)
	assert.Equal(t, 5, config.LimiterIPLimit)
	assert.Equal(t, 10, config.LimiterTokenLimit)
	assert.Equal(t, 300, config.BlockDuration)
}
