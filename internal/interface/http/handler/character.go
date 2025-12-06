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

// CharacterHandler serves the character details endpoint.
type CharacterHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewCharacterHandler creates a CharacterHandler configured with the v1 provider host.
func NewCharacterHandler(cacheManager *cache.Manager) *CharacterHandler {
	return &CharacterHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// GetCharacter returns character details with voice actors and animeography.
// @Summary      Get character details
// @Description  Returns detailed information about a specific character including voice actors and animeography
// @Tags         Characters
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Character ID"  example(asta-340)
// @Success      200  {object}  SuccessResponse  "Character details"
// @Failure      400  {object}  ErrorResponse    "Missing path parameter"
// @Failure      404  {object}  ErrorResponse    "Character not found"
// @Failure      502  {object}  ErrorResponse    "Failed to fetch character"
// @Router       /character/{id} [get]
func (h *CharacterHandler) GetCharacter(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	cacheKey := fmt.Sprintf("character:%s", id)
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryCharacter, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryCharacter,
		cacheKey,
		func() (interface{}, error) {
			character, err := extractors.ExtractCharacter(id, h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch character: %v", err))
			}

			if len(character.Results.Data) == 0 {
				return nil, fiber.NewError(fiber.StatusNotFound, "character not found")
			}
			return character, nil
		},
	)

	if err != nil {
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
