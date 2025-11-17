package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
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

	result, err := extractors.ExtractVoiceActorPage(id, page, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch character list: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": fiber.Map{
			"currentPage": page,
			"totalPages":  result.TotalPages,
			"data":        result.CharactersVoiceActors,
		},
	})
}
