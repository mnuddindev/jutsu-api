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

// NextEpisodeScheduleHandler serves the next episode schedule endpoint.
type NextEpisodeScheduleHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewNextEpisodeScheduleHandler creates a NextEpisodeScheduleHandler configured with the v1 provider host.
func NewNextEpisodeScheduleHandler(cacheManager *cache.Manager) *NextEpisodeScheduleHandler {
	return &NextEpisodeScheduleHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// GetNextEpisodeSchedule returns the next episode schedule for an anime.
// @Summary      Get schedule of next episode of anime
// @Description  Returns next-episode schedule information for a given anime
// @Tags         Schedule
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Anime ID or slug"  example(frieren-beyond-journeys-end-18542)
// @Success      200  {object}  SuccessResponse  "Next episode schedule"
// @Failure      400  {object}  ErrorResponse    "Missing path parameter"
// @Failure      502  {object}  ErrorResponse    "Failed to fetch next episode schedule"
// @Router       /schedule/{id} [get]
func (h *NextEpisodeScheduleHandler) GetNextEpisodeSchedule(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	cacheKey := fmt.Sprintf("next_episode_schedule:%s", id)
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryNextEpisode, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryNextEpisode,
		cacheKey,
		func() (interface{}, error) {
			nextSchedule, err := extractors.ExtractNextEpisodeSchedule(id, h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch next episode schedule: %v", err))
			}

			responseData := fiber.Map{
				"nextEpisodeSchedule": nextSchedule,
			}
			return responseData, nil
		},
	)

	if err != nil {
		// Handle Fiber errors
		if fErr, ok := err.(*fiber.Error); ok {
			return fErr
		}
		return fiber.NewError(500, err.Error())
	}

	var result fiber.Map
	json.Unmarshal(dataBytes, &result)

	return c.JSON(fiber.Map{
		"success": true,
		"cached":  wasCached,
		"results": result,
	})
}
