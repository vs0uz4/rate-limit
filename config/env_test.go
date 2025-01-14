package config

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	domainErrors "github.com/vs0uz4/rate-limit/internal/domain/errors"
)

func TestLoadConfigExecPathError(t *testing.T) {
	getExecutablePath = func() (string, error) {
		return "", errors.New("mock error")
	}
	defer func() {
		getExecutablePath = os.Executable
	}()

	config, err := LoadConfig()

	assert.Nil(t, config, "Configuration should be null when os.Executable fails")
	assert.Equal(t, domainErrors.ErrGettingExecPath, err, "Unexpected error returned")
}

func TestLoadConfigWithInvalidTokenLimits(t *testing.T) {
	envContent := `
REDIS_HOST=testhost
REDIS_PORT=1234
LIMITER_IP_LIMIT=10
BLOCK_DURATION=60
TOKEN_LIMITS=invalid_json
`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	config, err := LoadConfig()

	assert.Nil(t, config)
	assert.Error(t, err)
	assert.Equal(t, domainErrors.ErrInvalidTokenLimits, err, "Expected invalid token limits error")
}

func TestLoadConfig(t *testing.T) {
	envContent := `
REDIS_HOST=testhost
REDIS_PORT=1234
LIMITER_IP_LIMIT=10
BLOCK_DURATION=60
TOKEN_LIMITS={"token1":50,"token2":100}
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
	assert.Equal(t, 60, config.BlockDuration)
	assert.Equal(t, 50, config.TokenLimits["token1"])
	assert.Equal(t, 100, config.TokenLimits["token2"])
}

func TestLoadConfigWithEmptyTokenLimits(t *testing.T) {
	envContent := `
REDIS_HOST=testhost
REDIS_PORT=1234
LIMITER_IP_LIMIT=10
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
	assert.Equal(t, 60, config.BlockDuration)
	assert.Empty(t, config.TokenLimits)
}

func TestLoadConfigMissingEnvFile(t *testing.T) {
	os.Remove(".env")

	config, err := LoadConfig()

	assert.Nil(t, config, "Configuration should be null if .env is missing")
	assert.Error(t, err, "Expected an error when the .env file was missing")
	assert.Equal(t, domainErrors.ErrEnvFileNotFound, err, "Unexpected error message")
}

func TestLoadConfigMissingRedisHost(t *testing.T) {
	envContent := `
REDIS_PORT=1234
LIMITER_IP_LIMIT=10
BLOCK_DURATION=60
`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err, "Failed to create temporary .env file")
	defer os.Remove(".env")

	config, err := LoadConfig()

	assert.Nil(t, config, "Configuration should be null if REDIS_HOST is missing")
	assert.Error(t, err, "Expected an error when REDIS_HOST is missing")
	assert.Equal(t, domainErrors.ErrRedisHostRequired, err, "Unexpected error message")
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
	assert.Equal(t, 300, config.BlockDuration)
}
