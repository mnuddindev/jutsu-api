package router

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/internal/interface/http/handler"
	"github.com/mnuddindev/jutsu-api/internal/interface/middleware"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(app *fiber.App, cacheManager *cache.Manager) {
	// Health check routes
	healthHandler := handler.NewHealthHandler(cacheManager)
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)
	app.Get("/live", healthHandler.Live)

	// Swagger documentation JSON (served via swaggo)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Custom modern documentation shell
	swaggerUIDir := "./web/swagger"
	app.Static("/docs/assets", swaggerUIDir)
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(swaggerUIDir, "index.html"))
	})

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(swaggerUIDir, "favicon.ico"))
	})

	// app routes
	homeHandler := handler.NewHomeHandler(cacheManager)
	categoryHandler := handler.NewCategoryHandler(cacheManager)
	topTenHandler := handler.NewTopTenHandler(cacheManager)
	animeInfoHandler := handler.NewAnimeInfoHandler(cacheManager)
	episodeListHandler := handler.NewEpisodeListHandler(cacheManager)
	serversHandler := handler.NewServersHandler(cacheManager)
	streamHandler := handler.NewStreamHandler()
	searchHandler := handler.NewSearchHandler(cacheManager)
	suggestionHandler := handler.NewSuggestionHandler(cacheManager)
	filterHandler := handler.NewFilterHandler(cacheManager)
	scheduleHandler := handler.NewScheduleHandler(cacheManager)
	nextEpisodeScheduleHandler := handler.NewNextEpisodeScheduleHandler(cacheManager)
	randomHandler := handler.NewRandomHandler()
	qtipHandler := handler.NewQtipHandler(cacheManager)
	creatorHandler := handler.NewCreatorHandler(cacheManager)
	characterListHandler := handler.NewCharacterListHandler(cacheManager)
	characterHandler := handler.NewCharacterHandler(cacheManager)
	actorsHandler := handler.NewActorsHandler(cacheManager)
	watchlistHandler := handler.NewWatchlistHandler()
	topSearchHandler := handler.NewTopSearchHandler(cacheManager)

	api := app.Group("/api")

	// Home info routes
	api.Get("/", middleware.CacheByRoute(cacheManager, "home"), homeHandler.GetHomeInfo)

	// Category routes (genre and other categories)
	for _, routeType := range RouteTypes {
		api.Get("/"+routeType, categoryHandler.GetCategory)
	}

	// Top ten
	api.Get("/top-ten", topTenHandler.GetTopTen)

	// Anime info
	api.Get("/info", middleware.CacheByRoute(cacheManager, "anime"), animeInfoHandler.GetAnimeInfo)

	// Episodes
	api.Get("/episodes/:id", episodeListHandler.GetEpisodes)

	// Servers
	api.Get("/servers", serversHandler.GetServers)

	// Stream
	api.Get("/stream/:id", streamHandler.GetStream)
	api.Get("/stream/fallback/:id", streamHandler.GetStreamFallback)

	// Search
	api.Get("/search", searchHandler.Search)
	api.Get("/search/suggest", suggestionHandler.GetSuggestions)
	api.Get("/top-search", topSearchHandler.GetTopSearch)

	// Filter
	api.Get("/filter", filterHandler.Filter)

	// Schedule
	api.Get("/schedule", scheduleHandler.GetSchedule)
	api.Get("/schedule/:id", nextEpisodeScheduleHandler.GetNextEpisodeSchedule)

	// Random
	api.Get("/random", randomHandler.GetRandom)
	api.Get("/random/:id", randomHandler.GetRandomID)

	// Qtip
	api.Get("/qtip/:id", qtipHandler.GetQtip)

	// Producer
	api.Get("/producer/:id", creatorHandler.GetProducer)

	// Character and voice actors
	api.Get("/character/list/:id", characterListHandler.GetVoiceActors)
	api.Get("/character/:id", characterHandler.GetCharacter)
	api.Get("/actors/:id", actorsHandler.GetVoiceActor)

	// Watchlist
	api.Get("/watchlist/:userId", watchlistHandler.GetWatchlist)
	api.Get("/watchlist/:userId/:page", watchlistHandler.GetWatchlist)
}
