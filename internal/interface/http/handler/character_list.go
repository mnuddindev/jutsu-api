package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// CharacterListHandler serves the character list endpoint.
type CharacterListHandler struct {
	baseHost string
}

// NewCharacterListHandler creates a CharacterListHandler configured with the v1 provider host.
func NewCharacterListHandler() *CharacterListHandler {
	return &CharacterListHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetVoiceActors returns paginated list of characters with voice actors.
func (h *CharacterListHandler) GetVoiceActors(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	cacheKey := fmt.Sprintf("character_list:%s:%d", id, page)

	// Try to get from cache
	var cached map[string]interface{}
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	result, err := extractors.ExtractVoiceActorPage(id, page, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch character list: %v", err))
	}

	data := result.CharactersVoiceActors
	if data == nil {
		data = []extractors.CharactersVoiceActors{}
	}

	responseData := fiber.Map{
		"currentPage": page,
		"totalPages":  result.TotalPages,
		"data":        data,
	}

	// Cache the response
	_ = helper.SetCachedData(cacheKey, responseData, helper.CharacterCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": responseData,
	})
}
