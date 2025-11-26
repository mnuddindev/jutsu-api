package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// ServersHandler serves the available servers endpoint.
type ServersHandler struct {
	baseHost string
}

// NewServersHandler creates a ServersHandler configured with the v1 provider host.
func NewServersHandler() *ServersHandler {
	return &ServersHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetServers returns available streaming servers for an episode.
// @Summary      Get available servers of anime
// @Description  Returns the list of available streaming servers for an episode
// @Tags         Streaming
// @Accept       json
// @Produce      json
// @Param        ep   query     string  true  "Episode ID"  example(107257)
// @Success      200  {object}  SuccessResponse  "Servers list"
// @Failure      400  {object}  ErrorResponse    "Missing required query parameter"
// @Failure      502  {object}  ErrorResponse    "Failed to fetch servers"
// @Router       /servers [get]
func (h *ServersHandler) GetServers(c *fiber.Ctx) error {
	episodeID := strings.TrimSpace(c.Query("ep"))
	if episodeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'ep' is required")
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("servers:%s", episodeID)

	// Try to get from cache
	var cached []extractors.ServerItem
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && len(cached) > 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	servers, err := extractors.ExtractServers(episodeID, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch servers: %v", err))
	}

	// Return empty array instead of nil
	if servers == nil {
		servers = []extractors.ServerItem{}
	}

	// Cache the response
	_ = helper.SetCachedData(cacheKey, servers, helper.ServersCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": servers,
	})
}
