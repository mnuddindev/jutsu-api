package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// TopSearchHandler serves the top search endpoint.
type TopSearchHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewTopSearchHandler creates a TopSearchHandler configured with the v1 provider host.
func NewTopSearchHandler(cacheManager *cache.Manager) *TopSearchHandler {
	return &TopSearchHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// GetTopSearch returns top search keywords.
// @Summary      Get Top Search
// @Description  Returns the list of top search keywords
// @Tags         Search
// @Accept       json
// @Produce      json
// @Success      200  {object}  SuccessResponse  "Top search keywords"
// @Failure      502  {object}  ErrorResponse    "Failed to fetch top search"
// @Router       /top-search [get]
func (h *TopSearchHandler) GetTopSearch(c *fiber.Ctx) error {
	cacheKey := "top_search"
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryTopSearch, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryTopSearch,
		cacheKey,
		func() (interface{}, error) {
			topSearch, err := extractors.ExtractTopSearch(h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch top search: %v", err))
			}
			return topSearch, nil
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
