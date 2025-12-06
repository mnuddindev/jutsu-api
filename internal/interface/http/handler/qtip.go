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

// QtipHandler serves the qtip endpoint.
type QtipHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewQtipHandler creates a QtipHandler configured with the v1 provider host.
func NewQtipHandler(cacheManager *cache.Manager) *QtipHandler {
	return &QtipHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
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
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryQtip, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryQtip,
		cacheKey,
		func() (interface{}, error) {
			qtip, err := extractors.ExtractQtip(id, h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch qtip: %v", err))
			}
			return qtip, nil
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
