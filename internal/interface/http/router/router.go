package router

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"

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

	docsDir := "./web/docs"
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(docsDir, "index.html"))
	})

	app.Static("/docs/assets", filepath.Join(docsDir, "."))

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(docsDir, "favicon.ico"))
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

	api.Get("/",
		middleware.CacheByRoute(cacheManager, "home"),
		homeHandler.GetHomeInfo,
	)

	for _, routeType := range RouteTypes {
		api.Get("/"+routeType,
			middleware.CacheByRoute(cacheManager, "genre"),
			categoryHandler.GetCategory,
		)
	}

	api.Get("/top-ten",
		middleware.CacheByRoute(cacheManager, "top"),
		topTenHandler.GetTopTen,
	)

	api.Get("/info",
		middleware.CacheByRoute(cacheManager, "anime"),
		animeInfoHandler.GetAnimeInfo,
	)

	api.Get("/episodes/:id",
		middleware.CacheByRoute(cacheManager, "episodes"),
		episodeListHandler.GetEpisodes,
	)

	api.Get("/servers",
		middleware.CacheByRoute(cacheManager, "search"),
		serversHandler.GetServers,
	)

	api.Get("/search",
		middleware.CacheByRoute(cacheManager, "search"),
		searchHandler.Search,
	)
	api.Get("/search/suggest",
		middleware.CacheByRoute(cacheManager, "search"),
		suggestionHandler.GetSuggestions,
	)
	api.Get("/top-search",
		middleware.CacheByRoute(cacheManager, "top"),
		topSearchHandler.GetTopSearch,
	)

	api.Get("/filter",
		middleware.CacheByRoute(cacheManager, "search"),
		filterHandler.Filter,
	)

	api.Get("/schedule",
		middleware.CacheByRoute(cacheManager, "schedule"),
		scheduleHandler.GetSchedule,
	)
	api.Get("/schedule/:id",
		middleware.CacheByRoute(cacheManager, "schedule"),
		nextEpisodeScheduleHandler.GetNextEpisodeSchedule,
	)

	api.Get("/random", randomHandler.GetRandom)
	api.Get("/random/:id", randomHandler.GetRandomID)

	api.Get("/qtip/:id",
		middleware.CacheByRoute(cacheManager, "anime"),
		qtipHandler.GetQtip,
	)

	api.Get("/producer/:id",
		middleware.CacheByRoute(cacheManager, "genre"),
		creatorHandler.GetProducer,
	)

	api.Get("/character/list/:id",
		middleware.CacheByRoute(cacheManager, "character"),
		characterListHandler.GetVoiceActors,
	)
	api.Get("/character/:id",
		middleware.CacheByRoute(cacheManager, "character"),
		characterHandler.GetCharacter,
	)
	api.Get("/actors/:id",
		middleware.CacheByRoute(cacheManager, "character"),
		actorsHandler.GetVoiceActor,
	)

	api.Get("/stream/:id", streamHandler.GetStream)
	api.Get("/stream/fallback/:id", streamHandler.GetStreamFallback)

	api.Get("/watchlist/:userId", watchlistHandler.GetWatchlist)
	api.Get("/watchlist/:userId/:page", watchlistHandler.GetWatchlist)
}
