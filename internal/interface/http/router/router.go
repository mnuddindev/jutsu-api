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
	healthHandler := handler.NewHealthHandler(cacheManager)
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)
	app.Get("/live", healthHandler.Live)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	swaggerUIDir := "./web/swagger"
	app.Static("/docs/assets", swaggerUIDir)
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(swaggerUIDir, "index.html"))
	})

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(swaggerUIDir, "favicon.ico"))
	})

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

	api.Get("/", middleware.CacheByRoute(cacheManager, "home"), homeHandler.GetHomeInfo)

	for _, routeType := range RouteTypes {
		api.Get("/"+routeType, categoryHandler.GetCategory)
	}

	api.Get("/top-ten", topTenHandler.GetTopTen)

	api.Get("/info", middleware.CacheByRoute(cacheManager, "anime"), animeInfoHandler.GetAnimeInfo)

	api.Get("/episodes/:id", episodeListHandler.GetEpisodes)

	api.Get("/servers", serversHandler.GetServers)

	api.Get("/stream/:id", streamHandler.GetStream)
	api.Get("/stream/fallback/:id", streamHandler.GetStreamFallback)

	api.Get("/search", searchHandler.Search)
	api.Get("/search/suggest", suggestionHandler.GetSuggestions)
	api.Get("/top-search", topSearchHandler.GetTopSearch)

	api.Get("/filter", filterHandler.Filter)

	api.Get("/schedule", scheduleHandler.GetSchedule)
	api.Get("/schedule/:id", nextEpisodeScheduleHandler.GetNextEpisodeSchedule)

	api.Get("/random", randomHandler.GetRandom)
	api.Get("/random/:id", randomHandler.GetRandomID)

	api.Get("/qtip/:id", qtipHandler.GetQtip)

	api.Get("/producer/:id", creatorHandler.GetProducer)

	api.Get("/character/list/:id", characterListHandler.GetVoiceActors)
	api.Get("/character/:id", characterHandler.GetCharacter)
	api.Get("/actors/:id", actorsHandler.GetVoiceActor)

	api.Get("/watchlist/:userId", watchlistHandler.GetWatchlist)
	api.Get("/watchlist/:userId/:page", watchlistHandler.GetWatchlist)
}
