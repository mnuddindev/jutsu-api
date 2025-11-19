package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// CharacterHandler serves the character details endpoint.
type CharacterHandler struct {
	baseHost string
}

// NewCharacterHandler creates a CharacterHandler configured with the v1 provider host.
func NewCharacterHandler() *CharacterHandler {
	return &CharacterHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetCharacter returns character details with voice actors and animeography.
func (h *CharacterHandler) GetCharacter(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	cacheKey := fmt.Sprintf("character:%s", id)

	// Try to get from cache
	var cached interface{}
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	character, err := extractors.ExtractCharacter(id, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch character: %v", err))
	}

	// Check if data is empty
	if len(character.Results.Data) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "character not found")
	}

	// Cache the response
	_ = helper.SetCachedData(cacheKey, character.Results, helper.CharacterCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": character.Results,
	})
}
