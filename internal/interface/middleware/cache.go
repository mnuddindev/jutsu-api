package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
)

// CacheMiddleware provides caching functionality for responses using cache categories
// Use specific categories for better TTL management:
// - cache.CategoryHome for home endpoints
// - cache.CategorySearch for search endpoints
// - etc.
func CacheMiddleware(cacheManager *cache.Manager, category cache.CacheCategory) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only cache GET requests
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		// Skip if cache is disabled
		if !cacheManager.IsEnabled() {
			return c.Next()
		}

		ctx := c.Context()

		// Generate cache key from request
		cacheKey := generateCacheKey(c)

		// Try to get from cache using the manager
		cachedResponse, err := cacheManager.GetString(ctx, category, cacheKey)
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
				// Cache asynchronously to not block the response
				go func() {
					if err := cacheManager.SetString(ctx, category, cacheKey, string(responseBody)); err == nil {
						// Successfully cached
					}
				}()
				c.Set("X-Cache", "MISS")
			}
		}

		return nil
	}
}

// SimpleCacheMiddleware is a simpler version that uses a default category
// This is useful when you don't want to specify category for each route
func SimpleCacheMiddleware(cacheManager *cache.Manager) fiber.Handler {
	return CacheMiddleware(cacheManager, cache.CategoryHome)
}

// generateCacheKey generates a unique cache key for the request
func generateCacheKey(c *fiber.Ctx) string {
	// Include method, path, and query parameters
	key := fmt.Sprintf("%s:%s:%s", c.Method(), c.Path(), c.Queries())
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

// InvalidateCache invalidates cache for a specific pattern
func InvalidateCache(cacheManager *cache.Manager, category cache.CacheCategory, pattern string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Process request first
		if err := c.Next(); err != nil {
			return err
		}

		// Invalidate cache if request was successful
		if c.Response().StatusCode() < 400 {
			ctx := c.Context()
			_ = cacheManager.InvalidatePattern(ctx, category, pattern)
		}

		return nil
	}
}

// InvalidateCacheByPath invalidates cache based on the request path
// Useful for DELETE/PUT/POST operations that should clear related cache
func InvalidateCacheByPath(cacheManager *cache.Manager, category cache.CacheCategory) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Process request first
		if err := c.Next(); err != nil {
			return err
		}

		// Invalidate cache if request was successful
		if c.Response().StatusCode() < 400 {
			ctx := c.Context()
			// Generate pattern from path (e.g., /api/anime/123 -> *anime*123*)
			pattern := fmt.Sprintf("*%s*", c.Path())
			_ = cacheManager.InvalidatePattern(ctx, category, pattern)
		}

		return nil
	}
}

// CacheByRoute returns a cache middleware configured for specific route types
func CacheByRoute(cacheManager *cache.Manager, routeType string) fiber.Handler {
	// Map route types to cache categories
	categoryMap := map[string]cache.CacheCategory{
		"home":      cache.CategoryHome,
		"anime":     cache.CategoryAnimeInfo,
		"character": cache.CategoryCharacter,
		"episodes":  cache.CategoryEpisodes,
		"search":    cache.CategorySearch,
		"schedule":  cache.CategorySchedule,
		"genre":     cache.CategoryGenre,
		"top":       cache.CategoryTopTen,
		"stream":    cache.CategoryStream, // Usually 0 TTL
	}

	category, exists := categoryMap[routeType]
	if !exists {
		category = cache.CategoryHome // fallback
	}

	return CacheMiddleware(cacheManager, category)
}
