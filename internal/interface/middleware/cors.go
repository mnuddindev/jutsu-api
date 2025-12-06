package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/mnuddindev/jutsu-api/internal/config"
)

// SetupCORS configures CORS middleware
func SetupCORS(cfg *config.CorsConfig) fiber.Handler {
	methods := strings.Split(cfg.AllowMethods, ",")
	for i, method := range methods {
		methods[i] = strings.TrimSpace(method)
	}

	headers := strings.Split(cfg.AllowHeaders, ",")
	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
	}

	exposedHeaders := strings.Split(cfg.ExposeHeaders, ",")
	for i, header := range exposedHeaders {
		exposedHeaders[i] = strings.TrimSpace(header)
	}

	allowOrigins := strings.Join(cfg.AllowedOrigins, ",")
	if len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*" {
		if cfg.AllowCredentials {
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
