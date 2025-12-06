package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
)

type RateLimiterConfig struct {
	Max        int
	Expiration time.Duration
	RedisURL   string
}

// NewRateLimiter creates a rate limiter middleware
func NewRateLimiter(config RateLimiterConfig) fiber.Handler {
	var storage fiber.Storage

	if config.RedisURL != "" {
		storage = redis.New(redis.Config{
			URL:   config.RedisURL,
			Reset: false,
		})
	}

	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"message":     "Too many requests from this IP, please try again later",
				"retry_after": config.Expiration.Seconds(),
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.SlidingWindow{},
	})
}

// NewAPIKeyRateLimiter creates a rate limiter based on API key
func NewAPIKeyRateLimiter(config RateLimiterConfig) fiber.Handler {
	var storage fiber.Storage

	if config.RedisURL != "" {
		storage = redis.New(redis.Config{
			URL:   config.RedisURL,
			Reset: false,
		})
	}

	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			apiKey := c.Get("X-API-Key")
			if apiKey != "" {
				return "api_key:" + apiKey
			}
			return "ip:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"message":     "Too many requests, please try again later",
				"retry_after": config.Expiration.Seconds(),
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.SlidingWindow{},
	})
}

// RateLimitByEndpoint provides different rate limits for different endpoints
func RateLimitByEndpoint() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		var max int
		var expiration time.Duration

		switch {
		case path == "/api/" || path == "/api/home":
			max = 30
			expiration = 1 * time.Minute
		case path == "/api/search" || path == "/api/search/suggest":
			max = 60
			expiration = 1 * time.Minute
		case path == "/api/stream" || path == "/api/stream/fallback":
			max = 100
			expiration = 1 * time.Minute
		default:
			max = 100
			expiration = 1 * time.Minute
		}

		c.Locals("rate_limit_max", max)
		c.Locals("rate_limit_expiration", expiration)

		return c.Next()
	}
}
