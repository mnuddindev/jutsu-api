package handler

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// EpisodeListHandler serves the episode list endpoint.
type EpisodeListHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewEpisodeListHandler creates an EpisodeListHandler configured with the v1 provider host.
func NewEpisodeListHandler(cacheManager *cache.Manager) *EpisodeListHandler {
	return &EpisodeListHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// GetEpisodes returns the list of episodes for an anime.
// @Summary      Get anime episode list
// @Description  Returns the list of episodes for a specific anime
// @Tags         Anime
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Anime ID or slug"  example(frieren-beyond-journeys-end-18542)
// @Success      200  {object}  SuccessResponse  "Episode list"
// @Failure      400  {object}  ErrorResponse    "Missing or invalid path parameter"
// @Failure      502  {object}  ErrorResponse    "Failed to fetch episodes"
// @Router       /episodes/{id} [get]
func (h *EpisodeListHandler) GetEpisodes(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	encodedID := url.QueryEscape(id)
	ctx := c.Context()
	cacheKey := fmt.Sprintf("episodes_%s", encodedID)

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryEpisodes, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryEpisodes,
		cacheKey,
		func() (interface{}, error) {
			episodes, err := extractors.ExtractEpisodeList(encodedID, h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch episodes: %v", err))
			}
			return episodes, nil
		},
	)

	if err != nil {
		if fErr, ok := err.(*fiber.Error); ok {
			return fErr
		}
		return fiber.NewError(500, err.Error())
	}

	var result extractors.EpisodeList
	json.Unmarshal(dataBytes, &result)

	return c.JSON(fiber.Map{
		"success": true,
		"cached":  wasCached,
		"results": result,
	})
}
