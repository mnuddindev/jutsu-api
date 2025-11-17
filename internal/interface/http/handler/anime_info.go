package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// AnimeInfoHandler serves the anime info endpoint.
type AnimeInfoHandler struct {
	baseHost string
}

// NewAnimeInfoHandler creates an AnimeInfoHandler configured with the v1 provider host.
func NewAnimeInfoHandler() *AnimeInfoHandler {
	return &AnimeInfoHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetAnimeInfo returns combined anime info with seasons.
// @Summary      Get anime information
// @Description  Returns detailed information about a specific anime including seasons
// @Tags         Anime
// @Accept       json
// @Produce      json
// @Param        id   query     string  true  "Anime ID or slug"  example(frieren-beyond-journeys-end-18542)
// @Success      200  {object}  map[string]interface{}  "Anime information"
// @Failure      400  {object}  map[string]interface{}  "Bad Request"
// @Failure      502  {object}  map[string]interface{}  "Bad Gateway"
// @Router       /info [get]
func (h *AnimeInfoHandler) GetAnimeInfo(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Query("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'id' is required")
	}

	cacheKey := fmt.Sprintf("animeInfo_%s", id)
	var cached map[string]interface{}
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached != nil && len(cached) > 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	// Fetch anime info and seasons in parallel
	animeInfo, err1 := extractors.ExtractAnimeInfo(id, h.baseHost)
	seasons, err2 := extractors.ExtractSeasons(id, h.baseHost)

	if err1 != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch anime info: %v", err1))
	}
	if err2 != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch seasons: %v", err2))
	}

	responseData := map[string]interface{}{
		"data":    animeInfo,
		"seasons": seasons,
	}

	_ = helper.SetCachedData(cacheKey, responseData, helper.AnimeInfoCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": responseData,
	})
}
