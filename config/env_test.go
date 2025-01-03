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
	assert.NoError(t, err)
	defer os.Remove(".env")

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config)
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
	assert.Equal(t, errors.ErrMissingEnvFile, err, "Unexpected error message")
}
