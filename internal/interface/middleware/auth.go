package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/mnuddindev/jutsu-api/internal/repo"
	"github.com/mnuddindev/jutsu-api/internal/service"
)

// AuthMiddleware validates JWT access tokens
func AuthMiddleware(tokenService *service.TokenService, tokenRepo *repo.TokenRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format. Use: Bearer <token>",
			})
		}

		tokenString := parts[1]

		claims, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			if errors.Is(err, service.ErrExpiredToken) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token has expired",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or malformed token",
			})
		}

		isBlacklisted, err := tokenRepo.IsAccessTokenBlacklisted(c.Context(), claims.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify token",
			})
		}

		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token has been revoked",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("token_id", claims.ID)

		return c.Next()
	}
}

// OptionalAuthMiddleware extracts user info if token exists
func OptionalAuthMiddleware(tokenService *service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]
		claims, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Next()
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in context")
	}
	return userID, nil
}

// GetEmail extracts email from context
func GetEmail(c *fiber.Ctx) (string, error) {
	email, ok := c.Locals("email").(string)
	if !ok {
		return "", errors.New("email not found in context")
	}
	return email, nil
}

// GetRole extracts role from context
func GetRole(c *fiber.Ctx) (string, error) {
	role, ok := c.Locals("role").(string)
	if !ok {
		return "", errors.New("role not found in context")
	}
	return role, nil
}
