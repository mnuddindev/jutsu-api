package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// CharacterListHandler serves the character list endpoint.
type CharacterListHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewCharacterListHandler creates a CharacterListHandler configured with the v1 provider host.
func NewCharacterListHandler(cacheManager *cache.Manager) *CharacterListHandler {
	return &CharacterListHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// GetVoiceActors returns paginated list of characters with voice actors.
// @Summary      Get anime characters
// @Description  Returns a paginated list of characters with voice actors for a specific anime
// @Tags         Characters
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Anime ID or slug"  example(frieren-beyond-journeys-end-18542)
// @Param        page   query     int     false "Page number (must be >= 1)"  default(1)  example(1)
// @Success      200    {object}  SuccessResponse  "Characters and voice actors"
// @Failure      400    {object}  ErrorResponse    "Missing or invalid parameters"
// @Failure      502    {object}  ErrorResponse    "Failed to fetch character list"
// @Router       /character/list/{id} [get]
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
	ctx := c.Context()

	// Check if cached (for response field)
	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryCharacter, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryCharacter,
		cacheKey,
		func() (interface{}, error) {
			result, err := extractors.ExtractVoiceActorPage(id, page, h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch character list: %v", err))
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
			return responseData, nil
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
