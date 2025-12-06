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

// ActorsHandler serves the voice actor details endpoint.
type ActorsHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewActorsHandler creates an ActorsHandler configured with the v1 provider host.
func NewActorsHandler(cacheManager *cache.Manager) *ActorsHandler {
	return &ActorsHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
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
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryActor, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryActor,
		cacheKey,
		func() (interface{}, error) {
			voiceActor, err := extractors.ExtractVoiceActor(id, h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch voice actor: %v", err))
			}

			// Check if data is empty
			if len(voiceActor.Results.Data) == 0 {
				return nil, fiber.NewError(fiber.StatusNotFound, "no voice actor found")
			}
			return voiceActor, nil
		},
	)

	if err != nil {
		// Handle Fiber errors
		if fErr, ok := err.(*fiber.Error); ok {
			return fErr
		}
		return fiber.NewError(500, err.Error())
	}

	var result interface{}
	json.Unmarshal(dataBytes, &result)

	return c.JSON(fiber.Map{
		"success": true,
		"cached":  wasCached,
		"results": result,
	})
}
