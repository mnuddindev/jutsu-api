package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
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
func (h *ServersHandler) GetServers(c *fiber.Ctx) error {
	episodeID := strings.TrimSpace(c.Query("ep"))
	if episodeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'ep' is required")
	}

	servers, err := extractors.ExtractServers(episodeID, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch servers: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": servers,
	})
}
