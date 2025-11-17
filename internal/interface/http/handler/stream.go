package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// StreamHandler serves the stream and fallback endpoints.
type StreamHandler struct {
	baseHost string
}

// NewStreamHandler creates a StreamHandler that targets the v1 provider host.
func NewStreamHandler() *StreamHandler {
	return &StreamHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetStream resolves the default streaming info for an episode/server pair.
// @Summary      Get streaming information
// @Description  Returns streaming information for an episode with specified server and type
// @Tags         Streaming
// @Accept       json
// @Produce      json
// @Param        id       query     string  true  "Episode ID (format: anime-id?ep=episode-id)"  example(frieren-beyond-journeys-end-18542?ep=107257)
// @Param        server   query     string  false  "Server name (e.g., vidcloud, hd-1)"  example(hd-1)
// @Param        type     query     string  false  "Stream type (sub or dub)"  example(sub)
// @Success      200      {object}  map[string]interface{}  "Streaming information"
// @Failure      400      {object}  map[string]interface{}  "Bad Request"
// @Failure      502      {object}  map[string]interface{}  "Bad Gateway"
// @Router       /stream [get]
func (h *StreamHandler) GetStream(c *fiber.Ctx) error {
	return h.handleStreamRequest(c, false)
}

// GetStreamFallback resolves the fallback streaming info for an episode/server pair.
// @Summary      Get fallback streaming information
// @Description  Returns fallback streaming information when primary stream fails
// @Tags         Streaming
// @Accept       json
// @Produce      json
// @Param        id       query     string  true  "Episode ID (format: anime-id?ep=episode-id)"  example(frieren-beyond-journeys-end-18542?ep=107257)
// @Param        server   query     string  false  "Server name"  example(hd-1)
// @Param        type     query     string  false  "Stream type (sub or dub)"  example(sub)
// @Success      200      {object}  map[string]interface{}  "Fallback streaming information"
// @Failure      400      {object}  map[string]interface{}  "Bad Request"
// @Failure      503      {object}  map[string]interface{}  "Service Unavailable"
// @Router       /stream/fallback [get]
func (h *StreamHandler) GetStreamFallback(c *fiber.Ctx) error {
	return h.handleStreamRequest(c, true)
}

func (h *StreamHandler) handleStreamRequest(c *fiber.Ctx, fallback bool) error {
	episodeID := strings.TrimSpace(c.Query("id"))
	if episodeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'id' is required")
	}
	serverName := c.Query("server")
	streamType := c.Query("type")

	streamInfo, err := extractors.ExtractStreamingInfo(episodeID, serverName, streamType, fallback, h.baseHost)
	if err != nil {
		status := fiber.StatusBadGateway
		if fallback {
			status = fiber.StatusServiceUnavailable
		}
		return fiber.NewError(status, fmt.Sprintf("failed to resolve streaming info: %v", err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"results": streamInfo,
	})
}
