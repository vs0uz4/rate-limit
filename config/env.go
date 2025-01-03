package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/vs0uz4/rate-limit/internal/domain/errors"
)

type Config struct {
	RedisHost         string
	RedisPort         string
	LimiterIPLimit    int
	LimiterTokenLimit int
	BlockDuration     int
}

func LoadConfig() (*Config, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, errors.ErrGettingExecPath
	}

	basePath := filepath.Dir(ex)

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(basePath)
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("LIMITER_IP_LIMIT", 5)
	viper.SetDefault("LIMITER_TOKEN_LIMIT", 10)
	viper.SetDefault("BLOCK_DURATION", 300)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.ErrEnvFileNotFound
	}

	redisHost := viper.GetString("REDIS_HOST")
	if redisHost == "" {
		return nil, errors.ErrRedisHostRequired
	}

	return &Config{
		RedisHost:         viper.GetString("REDIS_HOST"),
		RedisPort:         viper.GetString("REDIS_PORT"),
		LimiterIPLimit:    viper.GetInt("LIMITER_IP_LIMIT"),
		LimiterTokenLimit: viper.GetInt("LIMITER_TOKEN_LIMIT"),
		BlockDuration:     viper.GetInt("BLOCK_DURATION"),
	}, nil
}
