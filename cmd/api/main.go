package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/mnuddindev/jutsu-api/internal/config"
	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	appLogger "github.com/mnuddindev/jutsu-api/internal/infrastructure/logger"
	"github.com/mnuddindev/jutsu-api/internal/interface/http/router"
	"github.com/mnuddindev/jutsu-api/internal/interface/middleware"
	"github.com/mnuddindev/jutsu-api/internal/interface/validation"
	"github.com/mnuddindev/jutsu-api/pkg/utils"

	_ "github.com/mnuddindev/jutsu-api/docs"
)

// @title           Jutsu API
// @version         1.0.0
// @description     High-Performance Anime Streaming API built with Go
// @description     A blazing-fast, production-ready RESTful API for anime streaming platforms
// @description
// @description     **Features:**
// @description     - Complete anime data (episodes, metadata, streaming sources)
// @description     - Advanced search and filtering
// @description     - Schedule management and notifications
// @description     - Character and voice actor information
// @description     - User watchlist support
// @description     - Optimized caching with Redis
// @description
// @description     **Disclaimer:**
// @description     This API does not store any files. It only links to media hosted on 3rd party services.
// @description     This API is explicitly made for educational purposes only.
//
// @contact.name   Jutsu API Support
// @contact.url    https://github.com/mnuddindev/jutsu-api
// @contact.email  support@jutsu-api.com
//
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
//
// @host      localhost:8080
// @BasePath  /api
//
// @schemes   http https
//
// @tag.name  Home
// @tag.description Home page data and featured content
//
// @tag.name  Anime
// @tag.description Anime information, episodes, and details
//
// @tag.name  Categories
// @tag.description Genre and category listings
//
// @tag.name  Search
// @tag.description Search and filter anime
//
// @tag.name  Streaming
// @tag.description Streaming information and servers
//
// @tag.name  Schedule
// @tag.description Anime schedule and episode notifications
//
// @tag.name  Characters
// @tag.description Character and voice actor information
//
// @tag.name  Random
// @tag.description Random anime discovery
//
// @tag.name  Watchlist
// @tag.description User watchlist management
//
// @tag.name  Health
// @tag.description Health check endpoints

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	if err := appLogger.InitLogger(&cfg.Logger); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer appLogger.Sync()

	utils.SetAppStartTime(time.Now())
	appLogger.Info("Starting Jutsu API", zap.String("version", cfg.App.Version))

	// Initialize validator
	if err := validation.InitValidator(); err != nil {
		appLogger.Fatal("Failed to initialize validator", zap.Error(err))
	}

	// Initialize cache (optional - app can run without it)
	var cacheManager *cache.Manager
	if err := cache.InitCache(&cfg.Redis); err != nil {
		appLogger.Warn("Failed to initialize cache - continuing without cache", zap.Error(err))
		// Create disabled cache manager
		cacheManager = cache.NewManager(nil, appLogger.Logger, cache.Config{Enabled: false})
	} else {
		appLogger.Info("Cache initialized successfully")

		// Get Redis client
		redisClient := cache.GetRedisClient()

		// Check if Redis is actually enabled by checking environment variable
		redisEnabled := cfg.Redis.Enabled

		// Create cache manager with TTL configuration from environment
		cacheConfig := cache.Config{
			Enabled:        redisEnabled,
			CharacterTTL:   cfg.Cache.CharacterTTL,   // 24 hours
			ActorTTL:       cfg.Cache.ActorTTL,       // 24 hours
			StudioTTL:      cfg.Cache.StudioTTL,      // 24 hours
			ProducerTTL:    cfg.Cache.ProducerTTL,    // 24 hours
			AnimeInfoTTL:   cfg.Cache.AnimeInfoTTL,   // 6 hours
			EpisodesTTL:    cfg.Cache.EpisodesTTL,    // 1 hour
			QtipTTL:        cfg.Cache.QtipTTL,        // 6 hours
			HomeTTL:        cfg.Cache.HomeTTL,        // 15 minutes
			TopTenTTL:      cfg.Cache.TopTenTTL,      // 1 hour
			TopSearchTTL:   cfg.Cache.TopSearchTTL,   // 1 hour
			TopAiringTTL:   cfg.Cache.TopAiringTTL,   // 2 hours
			GenreTTL:       cfg.Cache.GenreTTL,       // 4 hours
			ScheduleTTL:    cfg.Cache.ScheduleTTL,    // 10 minutes
			NextEpisodeTTL: cfg.Cache.NextEpisodeTTL, // 30 minutes
			ServersTTL:     cfg.Cache.ServersTTL,     // 5 minutes
			SearchTTL:      cfg.Cache.SearchTTL,      // 5 minutes
			SuggestTTL:     cfg.Cache.SuggestTTL,     // 3 minutes
			RandomTTL:      cfg.Cache.RandomTTL,      // 1 minute
			StreamTTL:      cfg.Cache.StreamTTL,      // Don't cache
		}

		cacheManager = cache.NewManager(redisClient, appLogger.Logger, cacheConfig)
		appLogger.Info("Cache manager initialized",
			zap.Bool("enabled", redisEnabled),
			zap.Int("anime_info_ttl", cacheConfig.AnimeInfoTTL))

		defer func() {
			if err := cache.CloseCache(); err != nil {
				appLogger.Error("Failed to close cache", zap.Error(err))
			}
		}()
	}

	// Create Fiber app with performance optimizations
	app := fiber.New(fiber.Config{
		AppName:                 cfg.App.Name,
		Prefork:                 cfg.Server.Prefork,
		ServerHeader:            "Jutsu",
		StrictRouting:           true,
		CaseSensitive:           true,
		DisableDefaultDate:      false,
		DisableStartupMessage:   false,
		ReadTimeout:             time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:            time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:             time.Duration(cfg.Server.IdleTimeout) * time.Second,
		EnablePrintRoutes:       cfg.App.Debug,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
		ProxyHeader:             fiber.HeaderXForwardedFor,
		ErrorHandler:            errorHandler,
	})

	// Setup middleware
	setupMiddleware(app, cfg)

	// Setup routes
	router.SetupRoutes(app, cacheManager)

	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine for graceful shutdown
	addr := cfg.GetServerAddr()
	go func() {
		if err := app.Listen(addr); err != nil {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Small delay to ensure Fiber startup message is printed
	time.Sleep(100 * time.Millisecond)
	appLogger.Info("Server started successfully", zap.String("address", addr))

	// Wait for interrupt signal to gracefully shutdown the server
	<-quit

	appLogger.Info("Shutting down server...")

	// Gracefully shutdown server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exited")
}

// setupMiddleware sets up all middleware
func setupMiddleware(app *fiber.App, cfg *config.Config) {
	// Recovery middleware (should be first)
	app.Use(middleware.SetupRecover())

	// CORS middleware
	app.Use(middleware.SetupCORS(&cfg.Cors))

	// Request logger middleware
	app.Use(middleware.RequestLogger())

	// Rate limiting
	if cfg.App.RateLimitEnabled {
		// Apply to all API routes
		app.Use("/api", middleware.NewRateLimiter(middleware.RateLimiterConfig{
			Max: 100,
		}))
		appLogger.Info("Rate limiting enabled",
			zap.Int("max", 100),
			zap.Int("window_seconds", 60))
	} else {
		appLogger.Info("Rate limiting disabled")
	}
}

// errorHandler is a custom error handler
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Skip logging for static assets like favicon.ico
	if c.Path() != "/favicon.ico" {
		appLogger.Error("Request Error",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Int("status", code),
		)
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}
