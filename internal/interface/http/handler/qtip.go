package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// QtipHandler serves the qtip endpoint.
type QtipHandler struct {
	baseHost string
}

// NewQtipHandler creates a QtipHandler configured with the v1 provider host.
func NewQtipHandler() *QtipHandler {
	return &QtipHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetQtip returns qtip data for an anime.
// @Summary      Get anime Qtip info
// @Description  Returns Qtip sidebar information for a specific anime
// @Tags         Anime
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Anime ID or slug"  example(frieren-beyond-journeys-end-18542)
// @Success      200  {object}  SuccessResponse  "Qtip info"
// @Failure      400  {object}  ErrorResponse    "Missing path parameter"
// @Failure      502  {object}  ErrorResponse    "Failed to fetch qtip info"
// @Router       /qtip/{id} [get]
func (h *QtipHandler) GetQtip(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	cacheKey := fmt.Sprintf("qtip:%s", id)

	// Try to get from cache
	var cached interface{}
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	qtip, err := extractors.ExtractQtip(id, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch qtip: %v", err))
	}

	// Cache the response
	_ = helper.SetCachedData(cacheKey, qtip, helper.QtipCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": qtip,
	})
}
