package config

import (
	"log"

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
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Error loading .env file:", err)
		return nil, errors.ErrMissingEnvFile
	}

	return &Config{
		RedisHost:         viper.GetString("REDIS_HOST"),
		RedisPort:         viper.GetString("REDIS_PORT"),
		LimiterIPLimit:    viper.GetInt("LIMITER_IP_LIMIT"),
		LimiterTokenLimit: viper.GetInt("LIMITER_TOKEN_LIMIT"),
		BlockDuration:     viper.GetInt("BLOCK_DURATION"),
	}, nil
}
