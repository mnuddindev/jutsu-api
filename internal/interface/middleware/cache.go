package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
)

// CacheMiddleware provides caching functionality for responses
func CacheMiddleware(duration time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only cache GET requests
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		// Generate cache key
		cacheKey := generateCacheKey(c)

		// Try to get from cache
		cachedResponse, err := cache.GetString(cacheKey)
		if err == nil && cachedResponse != "" {
			c.Set("Content-Type", "application/json")
			c.Set("X-Cache", "HIT")
			return c.SendString(cachedResponse)
		}

		// Process request
		if err := c.Next(); err != nil {
			return err
		}

		// Cache the response if status is 200
		if c.Response().StatusCode() == fiber.StatusOK {
			responseBody := c.Response().Body()
			if len(responseBody) > 0 {
				if err := cache.SetString(cacheKey, string(responseBody), duration); err == nil {
					c.Set("X-Cache", "MISS")
				}
			}
		}

		return nil
	}
}

// generateCacheKey generates a unique cache key for the request
func generateCacheKey(c *fiber.Ctx) string {
	key := fmt.Sprintf("%s:%s:%s", c.Method(), c.Path(), c.Queries())
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("cache:%s", hex.EncodeToString(hash[:]))
}

// InvalidateCache invalidates cache for a specific pattern
func InvalidateCache(pattern string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Process request first
		if err := c.Next(); err != nil {
			return err
		}

		// Invalidate cache if request was successful
		if c.Response().StatusCode() < 400 {
			_ = cache.DeletePattern(pattern)
		}

		return nil
	}
}
