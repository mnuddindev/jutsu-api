package handler

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// EpisodeListHandler serves the episode list endpoint.
type EpisodeListHandler struct {
	baseHost string
}

// NewEpisodeListHandler creates an EpisodeListHandler configured with the v1 provider host.
func NewEpisodeListHandler() *EpisodeListHandler {
	return &EpisodeListHandler{
		baseHost: utils.GetV1BaseHost(),
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

	// URL encode the id
	encodedID := url.QueryEscape(id)

	cacheKey := fmt.Sprintf("episodes_%s", encodedID)
	var cached extractors.EpisodeList
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && len(cached.Episodes) > 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	episodes, err := extractors.ExtractEpisodeList(encodedID, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch episodes: %v", err))
	}

	_ = helper.SetCachedData(cacheKey, episodes, helper.EpisodeListCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": episodes,
	})
}
