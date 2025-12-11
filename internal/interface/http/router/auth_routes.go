package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/internal/interface/http/handler"
	"github.com/mnuddindev/jutsu-api/internal/interface/middleware"
	"github.com/mnuddindev/jutsu-api/internal/repo"
	"github.com/mnuddindev/jutsu-api/internal/service"
)

// SetupAuthRoutes configures all authentication routes
func SetupAuthRoutes(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	tokenService *service.TokenService,
	tokenRepo *repo.TokenRepository,
	cacheManager *cache.Manager,
) {
	// Auth routes group
	user := app.Group("/api/user")

	// Public routes (
	// no authentication)
	user.Post("/register", authHandler.Register)
	user.Post("/login", authHandler.Login)
	user.Post("/refresh", authHandler.RefreshToken)

	// Protected routes (require authentication)
	protected := app.Group("/auth", middleware.AuthMiddleware(tokenService, tokenRepo))

	protected.Get("/profile", authHandler.GetProfile)
	protected.Patch("/profile", authHandler.UpdateProfile)
	protected.Post("/change-password", authHandler.ChangePassword)
	protected.Post("/logout", authHandler.Logout)
	protected.Post("/logout-all", authHandler.LogoutAll)

	// Admin only routes example
	admin := protected.Group("", middleware.RoleMiddleware("admin"))
	admin.Get("/admin/users", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Admin route - list users",
		})
	})
}
