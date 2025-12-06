package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/mnuddindev/jutsu-api/internal/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitCache initializes the Redis connection
func InitCache(cfg *config.RedisConfig) error {
	if !cfg.Enabled {
		return nil
	}

	fmt.Print(cfg.Port)
	redisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return redisClient
}

// CloseCache closes the Redis connection
func CloseCache() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
