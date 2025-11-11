package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/mnuddindev/jutsu-api/internal/config"
)

// SetupCORS configures CORS middleware
func SetupCORS(cfg *config.CorsConfig) fiber.Handler {
	// Parse allowed methods
	methods := strings.Split(cfg.AllowMethods, ",")
	for i, method := range methods {
		methods[i] = strings.TrimSpace(method)
	}

	// Parse allowed headers
	headers := strings.Split(cfg.AllowHeaders, ",")
	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
	}

	// Parse exposed headers
	exposedHeaders := strings.Split(cfg.ExposeHeaders, ",")
	for i, header := range exposedHeaders {
		exposedHeaders[i] = strings.TrimSpace(header)
	}

	// Handle wildcard origins
	allowOrigins := strings.Join(cfg.AllowedOrigins, ",")
	if len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*" {
		// If wildcard is used, credentials must be false
		if cfg.AllowCredentials {
			// Force credentials to false when using wildcard
			cfg.AllowCredentials = false
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     strings.Join(methods, ","),
		AllowHeaders:     strings.Join(headers, ","),
		AllowCredentials: cfg.AllowCredentials,
		ExposeHeaders:    strings.Join(exposedHeaders, ","),
		MaxAge:           cfg.MaxAge,
	})
}
