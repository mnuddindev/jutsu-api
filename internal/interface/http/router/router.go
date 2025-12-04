package router

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/mnuddindev/jutsu-api/internal/interface/http/handler"
	"github.com/mnuddindev/jutsu-api/internal/interface/middleware"
)

// SetupRoutes sets up all routes for the apilication
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", middleware.NewRateLimiter(middleware.RateLimiterConfig{}))

	// Health check routes
	healthHandler := handler.NewHealthHandler()
	api.Get("/health", healthHandler.Health)
	api.Get("/ready", healthHandler.Ready)
	api.Get("/live", healthHandler.Live)

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

	// API routes
	homeHandler := handler.NewHomeHandler()
	categoryHandler := handler.NewCategoryHandler()
	topTenHandler := handler.NewTopTenHandler()
	animeInfoHandler := handler.NewAnimeInfoHandler()
	episodeListHandler := handler.NewEpisodeListHandler()
	serversHandler := handler.NewServersHandler()
	streamHandler := handler.NewStreamHandler()
	searchHandler := handler.NewSearchHandler()
	suggestionHandler := handler.NewSuggestionHandler()
	filterHandler := handler.NewFilterHandler()
	scheduleHandler := handler.NewScheduleHandler()
	nextEpisodeScheduleHandler := handler.NewNextEpisodeScheduleHandler()
	randomHandler := handler.NewRandomHandler()
	qtipHandler := handler.NewQtipHandler()
	creatorHandler := handler.NewCreatorHandler()
	characterListHandler := handler.NewCharacterListHandler()
	characterHandler := handler.NewCharacterHandler()
	actorsHandler := handler.NewActorsHandler()
	watchlistHandler := handler.NewWatchlistHandler()
	topSearchHandler := handler.NewTopSearchHandler()

	// Home info routes
	api.Get("/", homeHandler.GetHomeInfo)

	// Category routes (genre and other categories)
	for _, routeType := range RouteTypes {
		api.Get("/api/"+routeType, categoryHandler.GetCategory)
	}

	// Top ten
	api.Get("/api/top-ten", topTenHandler.GetTopTen)

	// Anime info
	api.Get("/api/info", animeInfoHandler.GetAnimeInfo)

	// Episodes
	api.Get("/api/episodes/:id", episodeListHandler.GetEpisodes)

	// Servers
	api.Get("/api/servers", serversHandler.GetServers)

	// Stream
	api.Get("/api/stream/:id", streamHandler.GetStream)
	api.Get("/api/stream/fallback/:id", streamHandler.GetStreamFallback)

	// Search
	api.Get("/api/search", searchHandler.Search)
	api.Get("/api/search/suggest", suggestionHandler.GetSuggestions)
	api.Get("/api/top-search", topSearchHandler.GetTopSearch)

	// Filter
	api.Get("/api/filter", filterHandler.Filter)

	// Schedule
	api.Get("/api/schedule", scheduleHandler.GetSchedule)
	api.Get("/api/schedule/:id", nextEpisodeScheduleHandler.GetNextEpisodeSchedule)

	// Random
	api.Get("/api/random", randomHandler.GetRandom)
	api.Get("/api/random/:id", randomHandler.GetRandomID)

	// Qtip
	api.Get("/api/qtip/:id", qtipHandler.GetQtip)

	// Producer
	api.Get("/api/producer/:id", creatorHandler.GetProducer)

	// Character and voice actors
	api.Get("/api/character/list/:id", characterListHandler.GetVoiceActors)
	api.Get("/api/character/:id", characterHandler.GetCharacter)
	api.Get("/api/actors/:id", actorsHandler.GetVoiceActor)

	// Watchlist
	api.Get("/api/watchlist/:userId", watchlistHandler.GetWatchlist)
	api.Get("/api/watchlist/:userId/:page", watchlistHandler.GetWatchlist)
}
