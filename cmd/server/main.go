package main

import (
	"log"
	"time"

	"github.com/vs0uz4/rate-limit/config"
	"github.com/vs0uz4/rate-limit/internal/infra/redis"
	"github.com/vs0uz4/rate-limit/internal/rate_limiter"
	"github.com/vs0uz4/rate-limit/internal/webserver"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	redisClient, err := redis.NewRedisClient(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	redisAdapter := redis.NewRedisAdapter(redisClient.Client)

	limiter := rate_limiter.NewRateLimiter(rate_limiter.Config{
		Limit:         cfg.LimiterIPLimit,
		BlockDuration: time.Duration(cfg.BlockDuration) * time.Second,
		TokenLimits:   cfg.TokenLimits,
	}, redisAdapter)

	if err := webserver.Start(cfg, limiter); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
