package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/mnuddindev/jutsu-api/internal/config"
	appLogger "github.com/mnuddindev/jutsu-api/internal/infrastructure/logger"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

// InitCache initializes the Redis cache client
func InitCache(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		MaxRetries:      cfg.MaxRetries,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		DialTimeout:     time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(cfg.WriteTimeout) * time.Second,
		PoolTimeout:     time.Duration(cfg.PoolTimeout) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.IdleTimeout) * time.Minute,
	})

	// Test connection
	if err := Client.Ping(ctx).Err(); err != nil {
		appLogger.Warn("Failed to connect to Redis", zap.Error(err))
		// Don't return error - cache is optional, app can run without it
		return nil
	}

	appLogger.Info("Redis cache connection established successfully")
	return nil
}

// CloseCache closes the Redis cache connection
func CloseCache() error {
	if Client != nil {
		if err := Client.Close(); err != nil {
			appLogger.Error("Failed to close Redis connection", zap.Error(err))
			return err
		}
		appLogger.Info("Redis cache connection closed")
	}
	return nil
}

// HealthCheck checks if the cache is healthy
func HealthCheck() error {
	if Client == nil {
		return fmt.Errorf("cache client is nil")
	}

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("cache ping failed: %w", err)
	}

	return nil
}

// Set sets a key-value pair in the cache
func Set(key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("cache client is not initialized")
	}

	data, err := json.Marshal(value)
	if err != nil {
		appLogger.LogCacheOperation("SET", key, err)
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = Client.Set(ctx, key, data, expiration).Err()
	appLogger.LogCacheOperation("SET", key, err)
	return err
}

// Get retrieves a value from the cache
func Get(key string, dest interface{}) error {
	if Client == nil {
		return fmt.Errorf("cache client is not initialized")
	}

	data, err := Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			appLogger.LogCacheOperation("GET", key, fmt.Errorf("key not found"))
			return fmt.Errorf("key not found: %s", key)
		}
		appLogger.LogCacheOperation("GET", key, err)
		return err
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		appLogger.LogCacheOperation("GET", key, err)
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	appLogger.LogCacheOperation("GET", key, nil)
	return nil
}

// Delete deletes a key from the cache
func Delete(key string) error {
	if Client == nil {
		return fmt.Errorf("cache client is not initialized")
	}

	err := Client.Del(ctx, key).Err()
	appLogger.LogCacheOperation("DELETE", key, err)
	return err
}

// DeletePattern deletes all keys matching a pattern
func DeletePattern(pattern string) error {
	if Client == nil {
		return fmt.Errorf("cache client is not initialized")
	}

	iter := Client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		appLogger.LogCacheOperation("DELETE_PATTERN", pattern, err)
		return err
	}

	if len(keys) > 0 {
		err := Client.Del(ctx, keys...).Err()
		appLogger.LogCacheOperation("DELETE_PATTERN", pattern, err)
		return err
	}

	return nil
}

// Exists checks if a key exists in the cache
func Exists(key string) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("cache client is not initialized")
	}

	count, err := Client.Exists(ctx, key).Result()
	if err != nil {
		appLogger.LogCacheOperation("EXISTS", key, err)
		return false, err
	}

	return count > 0, nil
}

// SetNX sets a key-value pair only if the key does not exist
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("cache client is not initialized")
	}

	data, err := json.Marshal(value)
	if err != nil {
		appLogger.LogCacheOperation("SETNX", key, err)
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	result, err := Client.SetNX(ctx, key, data, expiration).Result()
	appLogger.LogCacheOperation("SETNX", key, err)
	return result, err
}

// Increment increments the value of a key by 1
func Increment(key string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("cache client is not initialized")
	}

	result, err := Client.Incr(ctx, key).Result()
	appLogger.LogCacheOperation("INCREMENT", key, err)
	return result, err
}

// IncrementBy increments the value of a key by the given amount
func IncrementBy(key string, value int64) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("cache client is not initialized")
	}

	result, err := Client.IncrBy(ctx, key, value).Result()
	appLogger.LogCacheOperation("INCREMENT_BY", key, err)
	return result, err
}

// Decrement decrements the value of a key by 1
func Decrement(key string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("cache client is not initialized")
	}

	result, err := Client.Decr(ctx, key).Result()
	appLogger.LogCacheOperation("DECREMENT", key, err)
	return result, err
}

// DecrementBy decrements the value of a key by the given amount
func DecrementBy(key string, value int64) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("cache client is not initialized")
	}

	result, err := Client.DecrBy(ctx, key, value).Result()
	appLogger.LogCacheOperation("DECREMENT_BY", key, err)
	return result, err
}

// Expire sets an expiration time on a key
func Expire(key string, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("cache client is not initialized")
	}

	err := Client.Expire(ctx, key, expiration).Err()
	appLogger.LogCacheOperation("EXPIRE", key, err)
	return err
}

// TTL returns the time to live of a key
func TTL(key string) (time.Duration, error) {
	if Client == nil {
		return 0, fmt.Errorf("cache client is not initialized")
	}

	result, err := Client.TTL(ctx, key).Result()
	appLogger.LogCacheOperation("TTL", key, err)
	return result, err
}

// GetString retrieves a string value from the cache
func GetString(key string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("cache client is not initialized")
	}

	result, err := Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			appLogger.LogCacheOperation("GET_STRING", key, fmt.Errorf("key not found"))
			return "", fmt.Errorf("key not found: %s", key)
		}
		appLogger.LogCacheOperation("GET_STRING", key, err)
		return "", err
	}

	appLogger.LogCacheOperation("GET_STRING", key, nil)
	return result, nil
}

// SetString sets a string value in the cache
func SetString(key string, value string, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("cache client is not initialized")
	}

	err := Client.Set(ctx, key, value, expiration).Err()
	appLogger.LogCacheOperation("SET_STRING", key, err)
	return err
}

// HSet sets a field in a hash
func HSet(key string, field string, value interface{}) error {
	if Client == nil {
		return fmt.Errorf("cache client is not initialized")
	}

	data, err := json.Marshal(value)
	if err != nil {
		appLogger.LogCacheOperation("HSET", key, err)
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = Client.HSet(ctx, key, field, data).Err()
	appLogger.LogCacheOperation("HSET", key, err)
	return err
}

// HGet retrieves a field from a hash
func HGet(key string, field string, dest interface{}) error {
	if Client == nil {
		return fmt.Errorf("cache client is not initialized")
	}

	data, err := Client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			appLogger.LogCacheOperation("HGET", key, fmt.Errorf("field not found"))
			return fmt.Errorf("field not found: %s", field)
		}
		appLogger.LogCacheOperation("HGET", key, err)
		return err
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		appLogger.LogCacheOperation("HGET", key, err)
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	appLogger.LogCacheOperation("HGET", key, nil)
	return nil
}

// HGetAll retrieves all fields from a hash
func HGetAll(key string) (map[string]string, error) {
	if Client == nil {
		return nil, fmt.Errorf("cache client is not initialized")
	}

	result, err := Client.HGetAll(ctx, key).Result()
	appLogger.LogCacheOperation("HGETALL", key, err)
	return result, err
}

// HDel deletes one or more fields from a hash
func HDel(key string, fields ...string) error {
	if Client == nil {
		return fmt.Errorf("cache client is not initialized")
	}

	err := Client.HDel(ctx, key, fields...).Err()
	appLogger.LogCacheOperation("HDEL", key, err)
	return err
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return Client
}
