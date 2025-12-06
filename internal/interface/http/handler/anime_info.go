package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// AnimeInfoHandler serves the anime info endpoint with optimized caching
type AnimeInfoHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewAnimeInfoHandler creates an AnimeInfoHandler with cache manager injection
func NewAnimeInfoHandler(cacheManager *cache.Manager) *AnimeInfoHandler {
	return &AnimeInfoHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// GetAnimeInfo returns combined anime info with seasons using the new cache manager
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

	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryAnimeInfo, id)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryAnimeInfo,
		id,
		func() (interface{}, error) {
			animeInfo, err1 := extractors.ExtractAnimeInfo(id, h.baseHost)
			seasons, err2 := extractors.ExtractSeasons(id, h.baseHost)

			if err1 != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway,
					fmt.Sprintf("failed to fetch anime info: %v", err1))
			}
			if err2 != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway,
					fmt.Sprintf("failed to fetch seasons: %v", err2))
			}

			return AnimeInfoResponse{
				Data:    animeInfo,
				Seasons: seasons,
			}, nil
		},
	)

	if err != nil {
		if fErr, ok := err.(*fiber.Error); ok {
			return fErr
		}
		return fiber.NewError(fiber.StatusInternalServerError,
			fmt.Sprintf("failed to process request: %v", err))
	}

	var responseData AnimeInfoResponse
	if err := json.Unmarshal(dataBytes, &responseData); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError,
			fmt.Sprintf("failed to parse response data: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"cached":  wasCached,
		"results": map[string]interface{}{
			"data":    responseData.Data,
			"seasons": responseData.Seasons,
		},
	})
}
