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

// ServersHandler serves the available servers endpoint.
type ServersHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewServersHandler creates a ServersHandler configured with the v1 provider host.
func NewServersHandler(cacheManager *cache.Manager) *ServersHandler {
	return &ServersHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
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
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryServers, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryServers,
		cacheKey,
		func() (interface{}, error) {

			servers, err := extractors.ExtractServers(episodeID, h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch servers: %v", err))
			}

			// Return empty array instead of nil
			if servers == nil {
				servers = []extractors.ServerItem{}
			}
			return servers, nil
		},
	)

	if err != nil {
		// Handle Fiber errors
		if fErr, ok := err.(*fiber.Error); ok {
			return fErr
		}
		return fiber.NewError(500, err.Error())
	}

	var result []extractors.ServerItem
	json.Unmarshal(dataBytes, &result)

	return c.JSON(fiber.Map{
		"success": true,
		"cached":  wasCached,
		"results": result,
	})
}
