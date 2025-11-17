package handlers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/mnuddindev/jutsu-api/internal/interface/http/handler"
)

func TestNewHandlers(t *testing.T) {
	t.Run("HomeHandler", func(t *testing.T) {
		h := handler.NewHomeHandler()
		assert.NotNil(t, h)
	})

	t.Run("CategoryHandler", func(t *testing.T) {
		h := handler.NewCategoryHandler()
		assert.NotNil(t, h)
	})

	t.Run("TopTenHandler", func(t *testing.T) {
		h := handler.NewTopTenHandler()
		assert.NotNil(t, h)
	})

	t.Run("AnimeInfoHandler", func(t *testing.T) {
		h := handler.NewAnimeInfoHandler()
		assert.NotNil(t, h)
	})

	t.Run("EpisodeListHandler", func(t *testing.T) {
		h := handler.NewEpisodeListHandler()
		assert.NotNil(t, h)
	})

	t.Run("ServersHandler", func(t *testing.T) {
		h := handler.NewServersHandler()
		assert.NotNil(t, h)
	})

	t.Run("SearchHandler", func(t *testing.T) {
		h := handler.NewSearchHandler()
		assert.NotNil(t, h)
	})

	t.Run("FilterHandler", func(t *testing.T) {
		h := handler.NewFilterHandler()
		assert.NotNil(t, h)
	})

	t.Run("SuggestionHandler", func(t *testing.T) {
		h := handler.NewSuggestionHandler()
		assert.NotNil(t, h)
	})

	t.Run("ScheduleHandler", func(t *testing.T) {
		h := handler.NewScheduleHandler()
		assert.NotNil(t, h)
	})

	t.Run("NextEpisodeScheduleHandler", func(t *testing.T) {
		h := handler.NewNextEpisodeScheduleHandler()
		assert.NotNil(t, h)
	})

	t.Run("RandomHandler", func(t *testing.T) {
		h := handler.NewRandomHandler()
		assert.NotNil(t, h)
	})

	t.Run("QtipHandler", func(t *testing.T) {
		h := handler.NewQtipHandler()
		assert.NotNil(t, h)
	})

	t.Run("CreatorHandler", func(t *testing.T) {
		h := handler.NewCreatorHandler()
		assert.NotNil(t, h)
	})

	t.Run("CharacterListHandler", func(t *testing.T) {
		h := handler.NewCharacterListHandler()
		assert.NotNil(t, h)
	})

	t.Run("CharacterHandler", func(t *testing.T) {
		h := handler.NewCharacterHandler()
		assert.NotNil(t, h)
	})

	t.Run("ActorsHandler", func(t *testing.T) {
		h := handler.NewActorsHandler()
		assert.NotNil(t, h)
	})

	t.Run("WatchlistHandler", func(t *testing.T) {
		h := handler.NewWatchlistHandler()
		assert.NotNil(t, h)
	})

	t.Run("TopSearchHandler", func(t *testing.T) {
		h := handler.NewTopSearchHandler()
		assert.NotNil(t, h)
	})
}

func TestHandlerValidation(t *testing.T) {
	app := fiber.New()

	t.Run("AnimeInfoHandler requires id parameter", func(t *testing.T) {
		h := handler.NewAnimeInfoHandler()
		app.Get("/test", h.GetAnimeInfo)

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("EpisodeListHandler requires id parameter", func(t *testing.T) {
		h := handler.NewEpisodeListHandler()
		app.Get("/test/:id", h.GetEpisodes)

		req := httptest.NewRequest("GET", "/test/", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServersHandler requires ep parameter", func(t *testing.T) {
		h := handler.NewServersHandler()
		app.Get("/test", h.GetServers)

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("SuggestionHandler requires keyword parameter", func(t *testing.T) {
		h := handler.NewSuggestionHandler()
		app.Get("/test", h.GetSuggestions)

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
