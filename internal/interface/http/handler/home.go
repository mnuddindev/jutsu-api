package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/scrape"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// HomeHandler serves the home info endpoint.
type HomeHandler struct {
	baseHost string
}

// NewHomeHandler creates a HomeHandler configured with the v1 provider host.
func NewHomeHandler() *HomeHandler {
	return &HomeHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetHomeInfo returns the combined home page data.
// @Summary      Get home page data
// @Description  Returns combined home page data including spotlights, trending, top 10, schedule, and category previews
// @Tags         Home
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Home page data"
// @Failure      502  {object}  map[string]interface{}  "Bad Gateway"
// @Router       / [get]
// @Router       /api [get]
func (h *HomeHandler) GetHomeInfo(c *fiber.Ctx) error {
	cacheKey := "homeInfo"
	var cached map[string]interface{}
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached != nil && len(cached) > 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	// Get current date in YYYY-MM-DD format
	now := time.Now()
	dateStr := now.Format("2006-01-02")

	// Extract genres from route types (first 41 are genres)
	// These match the Node.js routeTypes array
	routeTypes := []string{
		"genre/action", "genre/adventure", "genre/cars", "genre/comedy", "genre/dementia",
		"genre/demons", "genre/drama", "genre/ecchi", "genre/fantasy", "genre/game",
		"genre/harem", "genre/historical", "genre/horror", "genre/isekai", "genre/josei",
		"genre/kids", "genre/magic", "genre/martial-arts", "genre/mecha", "genre/military",
		"genre/music", "genre/mystery", "genre/parody", "genre/police", "genre/psychological",
		"genre/romance", "genre/samurai", "genre/school", "genre/sci-fi", "genre/seinen",
		"genre/shoujo", "genre/shoujo-ai", "genre/shounen", "genre/shounen-ai", "genre/slice-of-life",
		"genre/space", "genre/sports", "genre/super-power", "genre/supernatural", "genre/thriller",
		"genre/vampire",
	}
	genres := make([]string, 0, len(routeTypes))
	for _, rt := range routeTypes {
		genres = append(genres, strings.Replace(rt, "genre/", "", 1))
	}

	// Fetch all data in parallel
	spotlights, err1 := extractors.ExtractSpotlights(h.baseHost)
	trending, err2 := extractors.ExtractTrending(h.baseHost)
	topTen, err3 := extractors.ExtractTopTen(h.baseHost)
	schedule, err4 := extractors.ExtractSchedule(dateStr, -330, h.baseHost)

	// Extract category pages
	topAiringItems, _, err5 := helper.ExtractPage(1, "top-airing", h.baseHost)
	mostPopularItems, _, err6 := helper.ExtractPage(1, "most-popular", h.baseHost)
	mostFavoriteItems, _, err7 := helper.ExtractPage(1, "most-favorite", h.baseHost)
	latestCompletedItems, _, err8 := helper.ExtractPage(1, "completed", h.baseHost)
	latestEpisodeItems, _, err9 := helper.ExtractPage(1, "recently-updated", h.baseHost)
	topUpcomingItems, _, err10 := helper.ExtractPage(1, "top-upcoming", h.baseHost)
	recentlyAddedItems, _, err11 := helper.ExtractPage(1, "recently-added", h.baseHost)

	// Check for errors
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return fiber.NewError(fiber.StatusBadGateway, "failed to fetch home data")
	}
	if err5 != nil || err6 != nil || err7 != nil || err8 != nil || err9 != nil || err10 != nil || err11 != nil {
		return fiber.NewError(fiber.StatusBadGateway, "failed to fetch category data")
	}

	// Extract first item from category results
	var topAiringData, mostPopularData, mostFavoriteData, latestCompletedData, latestEpisodeData, topUpcomingData, recentlyAddedData interface{}
	if items, ok := topAiringItems.([]scrape.ExtractedItem); ok && len(items) > 0 {
		topAiringData = items[0]
	}
	if items, ok := mostPopularItems.([]scrape.ExtractedItem); ok && len(items) > 0 {
		mostPopularData = items[0]
	}
	if items, ok := mostFavoriteItems.([]scrape.ExtractedItem); ok && len(items) > 0 {
		mostFavoriteData = items[0]
	}
	if items, ok := latestCompletedItems.([]scrape.ExtractedItem); ok && len(items) > 0 {
		latestCompletedData = items[0]
	}
	if items, ok := latestEpisodeItems.([]scrape.ExtractedItem); ok && len(items) > 0 {
		latestEpisodeData = items[0]
	}
	if items, ok := topUpcomingItems.([]scrape.ExtractedItem); ok && len(items) > 0 {
		topUpcomingData = items[0]
	}
	if items, ok := recentlyAddedItems.([]scrape.ExtractedItem); ok && len(items) > 0 {
		recentlyAddedData = items[0]
	}

	responseData := map[string]interface{}{
		"spotlights":      spotlights,
		"trending":        trending,
		"topTen":          topTen,
		"today":           map[string]interface{}{"schedule": schedule},
		"topAiring":       topAiringData,
		"mostPopular":     mostPopularData,
		"mostFavorite":    mostFavoriteData,
		"latestCompleted": latestCompletedData,
		"latestEpisode":   latestEpisodeData,
		"topUpcoming":     topUpcomingData,
		"recentlyAdded":   recentlyAddedData,
		"genres":          genres,
	}

	_ = helper.SetCachedData(cacheKey, responseData, helper.HomeInfoCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": responseData,
	})
}
