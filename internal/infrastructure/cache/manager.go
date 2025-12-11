package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// CacheCategory defines cache categories with their TTLs
type CacheCategory string

const (
	// Static data - rarely changes
	CategoryCharacter CacheCategory = "character"
	CategoryActor     CacheCategory = "actor"
	CategoryStudio    CacheCategory = "studio"
	CategoryProducer  CacheCategory = "producer"

	// Semi-static data - changes daily/weekly
	CategoryAnimeInfo CacheCategory = "anime_info"
	CategoryEpisodes  CacheCategory = "episodes"
	CategoryQtip      CacheCategory = "qtip"

	// Dynamic data - changes frequently
	CategoryHome      CacheCategory = "home"
	CategoryTopTen    CacheCategory = "top_ten"
	CategoryTopSearch CacheCategory = "top_search"
	CategoryTopAiring CacheCategory = "top_airing"
	CategoryGenre     CacheCategory = "genre"

	// Time-sensitive data
	CategorySchedule    CacheCategory = "schedule"
	CategoryNextEpisode CacheCategory = "next_episode"
	CategoryServers     CacheCategory = "servers"

	// Short-lived or conditional
	CategorySearch  CacheCategory = "search"
	CategorySuggest CacheCategory = "suggest"
	CategoryRandom  CacheCategory = "random"
	CategoryStream  CacheCategory = "stream"

	CategoryUserSession  CacheCategory = "user_session"
	CategoryUserProfile  CacheCategory = "user_profile"
	CategoryRateLimit    CacheCategory = "rate_limit"
	CategoryJWTBlacklist CacheCategory = "jwt_blacklist"
	CategoryRefreshToken CacheCategory = "refresh_token"
)

// Manager handles caching operations with different TTLs per category
type Manager struct {
	client  *redis.Client
	logger  *zap.Logger
	enabled bool
	ttls    map[CacheCategory]time.Duration
}

// Config holds cache configuration
type Config struct {
	Enabled        bool
	CharacterTTL   int
	ActorTTL       int
	StudioTTL      int
	ProducerTTL    int
	AnimeInfoTTL   int
	EpisodesTTL    int
	QtipTTL        int
	HomeTTL        int
	TopTenTTL      int
	TopSearchTTL   int
	TopAiringTTL   int
	GenreTTL       int
	ScheduleTTL    int
	NextEpisodeTTL int
	ServersTTL     int
	SearchTTL      int
	SuggestTTL     int
	RandomTTL      int
	StreamTTL      int

	UserSessionTTL  int
	UserProfileTTL  int
	RateLimitTTL    int
	JWTBlacklistTTL int
	RefreshTokenTTL int
}

// NewManager creates a new cache manager
func NewManager(client *redis.Client, logger *zap.Logger, config Config) *Manager {
	ttls := map[CacheCategory]time.Duration{
		// Static data (24 hours default)
		CategoryCharacter: time.Duration(config.CharacterTTL) * time.Second,
		CategoryActor:     time.Duration(config.ActorTTL) * time.Second,
		CategoryStudio:    time.Duration(config.StudioTTL) * time.Second,
		CategoryProducer:  time.Duration(config.ProducerTTL) * time.Second,

		// Semi-static data
		CategoryAnimeInfo: time.Duration(config.AnimeInfoTTL) * time.Second,
		CategoryEpisodes:  time.Duration(config.EpisodesTTL) * time.Second,
		CategoryQtip:      time.Duration(config.QtipTTL) * time.Second,

		// Dynamic data
		CategoryHome:      time.Duration(config.HomeTTL) * time.Second,
		CategoryTopTen:    time.Duration(config.TopTenTTL) * time.Second,
		CategoryTopSearch: time.Duration(config.TopSearchTTL) * time.Second,
		CategoryTopAiring: time.Duration(config.TopAiringTTL) * time.Second,
		CategoryGenre:     time.Duration(config.GenreTTL) * time.Second,

		// Time-sensitive
		CategorySchedule:    time.Duration(config.ScheduleTTL) * time.Second,
		CategoryNextEpisode: time.Duration(config.NextEpisodeTTL) * time.Second,
		CategoryServers:     time.Duration(config.ServersTTL) * time.Second,

		// Short-lived
		CategorySearch:  time.Duration(config.SearchTTL) * time.Second,
		CategorySuggest: time.Duration(config.SuggestTTL) * time.Second,
		CategoryRandom:  time.Duration(config.RandomTTL) * time.Second,
		CategoryStream:  time.Duration(config.StreamTTL) * time.Second,

		// Auth Categories
		CategoryUserSession:  time.Duration(config.UserSessionTTL) * time.Second,
		CategoryUserProfile:  time.Duration(config.UserProfileTTL) * time.Second,
		CategoryRateLimit:    time.Duration(config.RateLimitTTL) * time.Second,
		CategoryJWTBlacklist: time.Duration(config.JWTBlacklistTTL) * time.Second,
		CategoryRefreshToken: time.Duration(config.RefreshTokenTTL) * time.Second,
	}

	return &Manager{
		client:  client,
		logger:  logger,
		enabled: config.Enabled,
		ttls:    ttls,
	}
}

// Get retrieves data from cache
func (m *Manager) Get(ctx context.Context, category CacheCategory, key string) ([]byte, error) {
	if !m.enabled || m.client == nil {
		return nil, fmt.Errorf("cache disabled")
	}

	fullKey := m.buildKey(category, key)

	data, err := m.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			m.logger.Debug("cache miss",
				zap.String("category", string(category)),
				zap.String("key", key),
			)
			return nil, fmt.Errorf("cache miss")
		}
		m.logger.Error("cache get error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, err
	}

	m.logger.Debug("cache hit",
		zap.String("category", string(category)),
		zap.String("key", key),
	)

	return data, nil
}

// Set stores data in cache with appropriate TTL
func (m *Manager) Set(ctx context.Context, category CacheCategory, key string, value interface{}) error {
	if !m.enabled || m.client == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		m.logger.Error("cache marshal error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	fullKey := m.buildKey(category, key)
	ttl := m.getTTL(category)

	if ttl == 0 {
		m.logger.Debug("skipping cache (TTL=0)",
			zap.String("category", string(category)),
			zap.String("key", key),
		)
		return nil
	}

	err = m.client.Set(ctx, fullKey, data, ttl).Err()
	if err != nil {
		m.logger.Error("cache set error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Duration("ttl", ttl),
			zap.Error(err),
		)
		return err
	}

	m.logger.Debug("cache set",
		zap.String("category", string(category)),
		zap.String("key", key),
		zap.Duration("ttl", ttl),
	)

	return nil
}

// Delete removes data from cache
func (m *Manager) Delete(ctx context.Context, category CacheCategory, key string) error {
	if !m.enabled || m.client == nil {
		return nil
	}

	fullKey := m.buildKey(category, key)

	err := m.client.Del(ctx, fullKey).Err()
	if err != nil {
		m.logger.Error("cache delete error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	m.logger.Debug("cache deleted",
		zap.String("category", string(category)),
		zap.String("key", key),
	)

	return nil
}

// GetOrSet gets from cache or sets if not exists (cache-aside pattern)
func (m *Manager) GetOrSet(ctx context.Context, category CacheCategory, key string, fetchFn func() (interface{}, error)) ([]byte, error) {
	data, err := m.Get(ctx, category, key)
	if err == nil {
		return data, nil
	}

	value, err := fetchFn()
	if err != nil {
		return nil, err
	}

	go func() {
		if err := m.Set(context.Background(), category, key, value); err != nil {
			m.logger.Error("async cache set failed",
				zap.String("category", string(category)),
				zap.String("key", key),
				zap.Error(err),
			)
		}
	}()

	return json.Marshal(value)
}

// InvalidatePattern deletes all keys matching a pattern
func (m *Manager) InvalidatePattern(ctx context.Context, category CacheCategory, pattern string) error {
	if !m.enabled || m.client == nil {
		return nil
	}

	fullPattern := m.buildKey(category, pattern)

	iter := m.client.Scan(ctx, 0, fullPattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := m.client.Del(ctx, iter.Val()).Err(); err != nil {
			m.logger.Error("cache pattern delete error",
				zap.String("key", iter.Val()),
				zap.Error(err),
			)
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	m.logger.Info("cache pattern invalidated",
		zap.String("category", string(category)),
		zap.String("pattern", pattern),
	)

	return nil
}

// GetStats returns cache statistics
func (m *Manager) GetStats(ctx context.Context) (map[string]interface{}, error) {
	if !m.enabled || m.client == nil {
		return map[string]interface{}{"enabled": false}, nil
	}

	info, err := m.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, err
	}

	dbSize, err := m.client.DBSize(ctx).Result()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"enabled": true,
		"dbsize":  dbSize,
		"info":    info,
	}, nil
}

// buildKey creates a full cache key with category prefix
func (m *Manager) buildKey(category CacheCategory, key string) string {
	return fmt.Sprintf("jutsu:%s:%s", category, key)
}

// getTTL returns TTL for a category
func (m *Manager) getTTL(category CacheCategory) time.Duration {
	if ttl, ok := m.ttls[category]; ok {
		return ttl
	}
	return 1 * time.Hour
}

// IsEnabled returns whether caching is enabled
func (m *Manager) IsEnabled() bool {
	return m.enabled
}

// ported from old cache.go for hash operations
// HSet sets a field in a hash, using the category key prefix.
func (m *Manager) HSet(ctx context.Context, category CacheCategory, key, field string, value interface{}) error {
	if !m.enabled || m.client == nil {
		return nil
	}

	fullKey := m.buildKey(category, key)
	data, err := json.Marshal(value)
	if err != nil {
		m.logger.Error("cache HSet marshal error", zap.String("category", string(category)), zap.String("key", key), zap.Error(err))
		return err
	}

	err = m.client.HSet(ctx, fullKey, field, data).Err()
	if err != nil {
		m.logger.Error("cache HSet error", zap.String("category", string(category)), zap.String("key", key), zap.String("field", field), zap.Error(err))
		return err
	}
	m.logger.Debug("cache HSet", zap.String("category", string(category)), zap.String("key", key), zap.String("field", field))
	return nil
}

// HGet retrieves a field from a hash and unmarshals it into dest.
func (m *Manager) HGet(ctx context.Context, category CacheCategory, key, field string, dest interface{}) error {
	if !m.enabled || m.client == nil {
		return fmt.Errorf("cache disabled")
	}

	fullKey := m.buildKey(category, key)
	data, err := m.client.HGet(ctx, fullKey, field).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("field not found")
		}
		m.logger.Error("cache HGet error", zap.String("category", string(category)), zap.String("key", key), zap.String("field", field), zap.Error(err))
		return err
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		m.logger.Error("cache HGet unmarshal error", zap.String("category", string(category)), zap.String("key", key), zap.String("field", field), zap.Error(err))
		return err
	}
	m.logger.Debug("cache HGet hit", zap.String("category", string(category)), zap.String("key", key), zap.String("field", field))
	return nil
}

// HDel deletes one or more fields from a hash.
func (m *Manager) HDel(ctx context.Context, category CacheCategory, key string, fields ...string) error {
	if !m.enabled || m.client == nil {
		return nil
	}
	fullKey := m.buildKey(category, key)
	err := m.client.HDel(ctx, fullKey, fields...).Err()
	if err != nil {
		m.logger.Error("cache HDel error", zap.String("category", string(category)), zap.String("key", key), zap.Error(err))
		return err
	}
	m.logger.Debug("cache HDel", zap.String("category", string(category)), zap.String("key", key))
	return nil
}

// Increment increments the value of a key by 1.
func (m *Manager) Increment(ctx context.Context, category CacheCategory, key string) (int64, error) {
	if !m.enabled || m.client == nil {
		return 0, fmt.Errorf("cache disabled")
	}

	fullKey := m.buildKey(category, key)
	result, err := m.client.Incr(ctx, fullKey).Result()
	if err != nil {
		m.logger.Error("cache Incr error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, err
	}
	return result, nil
}

// IncrementBy increments the value of a key by the given amount.
func (m *Manager) IncrementBy(ctx context.Context, category CacheCategory, key string, value int64) (int64, error) {
	if !m.enabled || m.client == nil {
		return 0, fmt.Errorf("cache disabled")
	}

	fullKey := m.buildKey(category, key)
	result, err := m.client.IncrBy(ctx, fullKey, value).Result()
	if err != nil {
		m.logger.Error("cache IncrBy error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Int64("value", value),
			zap.Error(err),
		)
		return 0, err
	}
	return result, nil
}

// Decrement decrements the value of a key by 1.
func (m *Manager) Decrement(ctx context.Context, category CacheCategory, key string) (int64, error) {
	if !m.enabled || m.client == nil {
		return 0, fmt.Errorf("cache disabled")
	}

	fullKey := m.buildKey(category, key)
	result, err := m.client.Decr(ctx, fullKey).Result()
	if err != nil {
		m.logger.Error("cache Decr error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, err
	}
	return result, nil
}

// DecrementBy decrements the value of a key by the given amount.
func (m *Manager) DecrementBy(ctx context.Context, category CacheCategory, key string, value int64) (int64, error) {
	if !m.enabled || m.client == nil {
		return 0, fmt.Errorf("cache disabled")
	}

	fullKey := m.buildKey(category, key)
	result, err := m.client.DecrBy(ctx, fullKey, value).Result()
	if err != nil {
		m.logger.Error("cache DecrBy error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Int64("value", value),
			zap.Error(err),
		)
		return 0, err
	}
	return result, nil
}

// Exists checks if a key exists in the cache.
func (m *Manager) Exists(ctx context.Context, category CacheCategory, key string) (bool, error) {
	if !m.enabled || m.client == nil {
		return false, nil // Assume not exists if disabled
	}

	fullKey := m.buildKey(category, key)
	count, err := m.client.Exists(ctx, fullKey).Result()
	if err != nil {
		m.logger.Error("cache Exists error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Error(err),
		)
		return false, err
	}
	return count > 0, nil
}

// Expire sets an expiration time on a key.
func (m *Manager) Expire(ctx context.Context, category CacheCategory, key string, expiration time.Duration) error {
	if !m.enabled || m.client == nil {
		return nil
	}

	fullKey := m.buildKey(category, key)
	err := m.client.Expire(ctx, fullKey, expiration).Err()
	if err != nil {
		m.logger.Error("cache Expire error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Duration("expiration", expiration),
			zap.Error(err),
		)
		return err
	}
	m.logger.Debug("cache Expire set",
		zap.String("category", string(category)),
		zap.String("key", key),
		zap.Duration("expiration", expiration),
	)
	return nil
}

// TTL returns the remaining time to live of a key.
func (m *Manager) TTL(ctx context.Context, category CacheCategory, key string) (time.Duration, error) {
	if !m.enabled || m.client == nil {
		return 0, fmt.Errorf("cache disabled")
	}

	fullKey := m.buildKey(category, key)
	result, err := m.client.TTL(ctx, fullKey).Result()
	if err != nil {
		m.logger.Error("cache TTL error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, err
	}
	return result, nil
}

// SetNX sets a key-value pair only if the key does not exist.
// This uses the category's configured TTL.
func (m *Manager) SetNX(ctx context.Context, category CacheCategory, key string, value interface{}) (bool, error) {
	if !m.enabled || m.client == nil {
		return false, nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		m.logger.Error("cache SetNX marshal error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Error(err),
		)
		return false, err
	}

	fullKey := m.buildKey(category, key)
	ttl := m.getTTL(category)

	if ttl == 0 {
		return false, nil
	}

	result, err := m.client.SetNX(ctx, fullKey, data, ttl).Result()
	if err != nil {
		m.logger.Error("cache SetNX error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Duration("ttl", ttl),
			zap.Error(err),
		)
		return false, err
	}

	if result {
		m.logger.Debug("cache SetNX success",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Duration("ttl", ttl),
		)
	} else {
		m.logger.Debug("cache SetNX failed (key existed)",
			zap.String("category", string(category)),
			zap.String("key", key),
		)
	}

	return result, nil
}

// SetString sets a string value in the cache, applying the category's TTL.
func (m *Manager) SetString(ctx context.Context, category CacheCategory, key string, value string) error {
	if !m.enabled || m.client == nil {
		return nil
	}

	fullKey := m.buildKey(category, key)
	ttl := m.getTTL(category)

	// Don't cache if TTL is 0
	if ttl == 0 {
		m.logger.Debug("skipping string cache (TTL=0)",
			zap.String("category", string(category)),
			zap.String("key", key),
		)
		return nil
	}

	err := m.client.Set(ctx, fullKey, value, ttl).Err()
	if err != nil {
		m.logger.Error("cache SetString error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Duration("ttl", ttl),
			zap.Error(err),
		)
		return err
	}

	m.logger.Debug("cache SetString set",
		zap.String("category", string(category)),
		zap.String("key", key),
		zap.Duration("ttl", ttl),
	)

	return nil
}

// GetString retrieves a string value from the cache.
func (m *Manager) GetString(ctx context.Context, category CacheCategory, key string) (string, error) {
	if !m.enabled || m.client == nil {
		return "", fmt.Errorf("cache disabled")
	}

	fullKey := m.buildKey(category, key)

	result, err := m.client.Get(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			m.logger.Debug("cache string miss",
				zap.String("category", string(category)),
				zap.String("key", key),
			)
			return "", fmt.Errorf("key not found: %s", key)
		}
		m.logger.Error("cache GetString error",
			zap.String("category", string(category)),
			zap.String("key", key),
			zap.Error(err),
		)
		return "", err
	}

	m.logger.Debug("cache string hit",
		zap.String("category", string(category)),
		zap.String("key", key),
	)
	return result, nil
}

// Ping checks the connection to Redis by executing the PING command.
// This is the correct, port-ready replacement for the old global cache.HealthCheck.
func (m *Manager) Ping(ctx context.Context) error {
	if !m.enabled || m.client == nil {
		m.logger.Error("cache ping failed: client not initialized or disabled")
		return fmt.Errorf("cache client is not initialized or disabled")
	}

	if err := m.client.Ping(ctx).Err(); err != nil {
		m.logger.Error("cache ping failed", zap.Error(err))
		return fmt.Errorf("cache ping failed: %w", err)
	}

	m.logger.Debug("cache ping successful")
	return nil
}
