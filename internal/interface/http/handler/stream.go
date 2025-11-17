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
func (h *StreamHandler) GetStream(c *fiber.Ctx) error {
	return h.handleStreamRequest(c, false)
}

// GetStreamFallback resolves the fallback streaming info for an episode/server pair.
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
