package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// NextEpisodeScheduleHandler serves the next episode schedule endpoint.
type NextEpisodeScheduleHandler struct {
	baseHost string
}

// NewNextEpisodeScheduleHandler creates a NextEpisodeScheduleHandler configured with the v1 provider host.
func NewNextEpisodeScheduleHandler() *NextEpisodeScheduleHandler {
	return &NextEpisodeScheduleHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetNextEpisodeSchedule returns the next episode schedule for an anime.
func (h *NextEpisodeScheduleHandler) GetNextEpisodeSchedule(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	nextSchedule, err := extractors.ExtractNextEpisodeSchedule(id, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch next episode schedule: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": fiber.Map{
			"nextEpisodeSchedule": nextSchedule,
		},
	})
}
