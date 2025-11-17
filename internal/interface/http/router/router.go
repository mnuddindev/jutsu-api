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
	app.Get("/api", homeHandler.GetHomeInfo)
	app.Get("/api/", homeHandler.GetHomeInfo)

	// Category routes (genre and other categories)
	for _, routeType := range RouteTypes {
		app.Get("/api/"+routeType, categoryHandler.GetCategory)
	}

	// Top ten
	app.Get("/api/top-ten", topTenHandler.GetTopTen)

	// Anime info
	app.Get("/api/info", animeInfoHandler.GetAnimeInfo)

	// Episodes
	app.Get("/api/episodes/:id", episodeListHandler.GetEpisodes)

	// Servers
	app.Get("/api/servers/:id", serversHandler.GetServers)

	// Stream
	app.Get("/api/stream", streamHandler.GetStream)
	app.Get("/api/stream/fallback", streamHandler.GetStreamFallback)

	// Search
	app.Get("/api/search", searchHandler.Search)
	app.Get("/api/search/suggest", suggestionHandler.GetSuggestions)
	app.Get("/api/top-search", topSearchHandler.GetTopSearch)

	// Filter
	app.Get("/api/filter", filterHandler.Filter)

	// Schedule
	app.Get("/api/schedule", scheduleHandler.GetSchedule)
	app.Get("/api/schedule/:id", nextEpisodeScheduleHandler.GetNextEpisodeSchedule)

	// Random
	app.Get("/api/random", randomHandler.GetRandom)
	app.Get("/api/random/id", randomHandler.GetRandomID)

	// Qtip
	app.Get("/api/qtip/:id", qtipHandler.GetQtip)

	// Producer
	app.Get("/api/producer/:id", creatorHandler.GetProducer)
	app.Get("/api/studio/:id", creatorHandler.GetStudio)

	// Character and voice actors
	app.Get("/api/character/list/:id", characterListHandler.GetVoiceActors)
	app.Get("/api/character/:id", characterHandler.GetCharacter)
	app.Get("/api/actors/:id", actorsHandler.GetVoiceActor)

	// Watchlist
	app.Get("/api/watchlist/:userId", watchlistHandler.GetWatchlist)
	app.Get("/api/watchlist/:userId/:page", watchlistHandler.GetWatchlist)
}
