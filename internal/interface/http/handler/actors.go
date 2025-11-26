package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
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
// @Summary      Get voice actor details
// @Description  Returns detailed information about a specific voice actor and their roles
// @Tags         Characters
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Voice actor ID"  example(gakuto-kajiwara-534)
// @Success      200  {object}  SuccessResponse  "Voice actor details"
// @Failure      400  {object}  ErrorResponse    "Missing path parameter"
// @Failure      404  {object}  ErrorResponse    "No voice actor found"
// @Failure      502  {object}  ErrorResponse    "Failed to fetch voice actor"
// @Router       /actors/{id} [get]
func (h *ActorsHandler) GetVoiceActor(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	cacheKey := fmt.Sprintf("voice_actor:%s", id)

	// Try to get from cache
	var cached interface{}
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	voiceActor, err := extractors.ExtractVoiceActor(id, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch voice actor: %v", err))
	}

	// Check if data is empty
	if len(voiceActor.Results.Data) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no voice actor found")
	}

	// Cache the response
	_ = helper.SetCachedData(cacheKey, voiceActor.Results, helper.VoiceActorCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": voiceActor.Results,
	})
}
