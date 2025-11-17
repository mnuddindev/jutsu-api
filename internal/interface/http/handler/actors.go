package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// ActorsHandler serves the voice actor details endpoint.
type ActorsHandler struct {
	baseHost string
}

// NewActorsHandler creates an ActorsHandler configured with the v1 provider host.
func NewActorsHandler() *ActorsHandler {
	return &ActorsHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetVoiceActor returns voice actor details with roles.
func (h *ActorsHandler) GetVoiceActor(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	voiceActor, err := extractors.ExtractVoiceActor(id, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch voice actor: %v", err))
	}

	// Check if data is empty
	if len(voiceActor.Results.Data) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no voice actor found")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": voiceActor.Results,
	})
}
