package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/vs0uz4/rate-limit/internal/domain/errors"
)

var getExecutablePath = os.Executable

type Config struct {
	WebServerPort  string
	RedisHost      string
	RedisPort      string
	LimiterIPLimit int
	BlockDuration  int
	TokenLimits    map[string]int
}

func LoadConfig() (*Config, error) {
	ex, err := getExecutablePath()
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
	viper.SetDefault("BLOCK_DURATION", 300)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.ErrEnvFileNotFound
	}

	redisHost := viper.GetString("REDIS_HOST")
	if redisHost == "" {
		return nil, errors.ErrRedisHostRequired
	}

	tokenLimits := make(map[string]int)
	tokenLimitsJSON := viper.GetString("TOKEN_LIMITS")
	if tokenLimitsJSON != "" {
		err := json.Unmarshal([]byte(tokenLimitsJSON), &tokenLimits)
		if err != nil {
			return nil, errors.ErrInvalidTokenLimits
		}
	}

	return &Config{
		WebServerPort:  "8080",
		RedisHost:      viper.GetString("REDIS_HOST"),
		RedisPort:      viper.GetString("REDIS_PORT"),
		LimiterIPLimit: viper.GetInt("LIMITER_IP_LIMIT"),
		BlockDuration:  viper.GetInt("BLOCK_DURATION"),
		TokenLimits:    tokenLimits,
	}, nil
}
