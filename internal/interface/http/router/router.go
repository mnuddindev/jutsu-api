package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/fiber-swagger"

	"github.com/mnuddindev/jutsu-api/internal/interface/http/handler"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(app *fiber.App) {
	// Health check routes
	healthHandler := handler.NewHealthHandler()
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)
	app.Get("/live", healthHandler.Live)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// API routes (matching the Node.js API structure)
	notImplementedHandler := handler.NewNotImplementedHandler()
	streamHandler := handler.NewStreamHandler()
	creatorHandler := handler.NewCreatorHandler()
	randomHandler := handler.NewRandomHandler()

	// Home info routes
	app.Get("/api", notImplementedHandler.NotImplemented)
	app.Get("/api/", notImplementedHandler.NotImplemented)

	// Category routes (genre and other categories)
	for _, routeType := range RouteTypes {
		app.Get("/api/"+routeType, notImplementedHandler.NotImplemented)
	}

	// Top ten
	app.Get("/api/top-ten", notImplementedHandler.NotImplemented)

	// Anime info
	app.Get("/api/info", notImplementedHandler.NotImplemented)

	// Episodes
	app.Get("/api/episodes/:id", notImplementedHandler.NotImplemented)

	// Servers
	app.Get("/api/servers/:id", notImplementedHandler.NotImplemented)

	// Stream
	app.Get("/api/stream", streamHandler.GetStream)
	app.Get("/api/stream/fallback", streamHandler.GetStreamFallback)

	// Search
	app.Get("/api/search", notImplementedHandler.NotImplemented)
	app.Get("/api/search/suggest", notImplementedHandler.NotImplemented)
	app.Get("/api/top-search", notImplementedHandler.NotImplemented)

	// Filter
	app.Get("/api/filter", notImplementedHandler.NotImplemented)

	// Schedule
	app.Get("/api/schedule", notImplementedHandler.NotImplemented)
	app.Get("/api/schedule/:id", notImplementedHandler.NotImplemented)

	// Random
	app.Get("/api/random", randomHandler.GetRandom)
	app.Get("/api/random/id", randomHandler.GetRandomID)

	// Qtip
	app.Get("/api/qtip/:id", notImplementedHandler.NotImplemented)

	// Producer
	app.Get("/api/producer/:id", creatorHandler.GetProducer)
	app.Get("/api/studio/:id", creatorHandler.GetStudio)

	// Character and voice actors
	app.Get("/api/character/list/:id", notImplementedHandler.NotImplemented)
	app.Get("/api/character/:id", notImplementedHandler.NotImplemented)
	app.Get("/api/actors/:id", notImplementedHandler.NotImplemented)

	// Watchlist
	app.Get("/api/watchlist/:userId", notImplementedHandler.NotImplemented)
	app.Get("/api/watchlist/:userId/:page", notImplementedHandler.NotImplemented)
}
