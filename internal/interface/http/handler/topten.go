package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// TopTenHandler serves the top ten anime endpoint.
type TopTenHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewTopTenHandler creates a TopTenHandler configured with the v1 provider host.
func NewTopTenHandler(cacheManaer *cache.Manager) *TopTenHandler {
	return &TopTenHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManaer,
	}
}

// GetTopTen returns the top 10 anime data.
// @Summary      Get Top 10 anime's info
// @Description  Returns the Top 10 anime lists for today, week and month
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  SuccessResponse  "Top 10 anime data"
// @Failure      502  {object}  ErrorResponse    "Failed to fetch top ten"
// @Router       /top-ten [get]
func (h *TopTenHandler) GetTopTen(c *fiber.Ctx) error {
	cacheKey := "topTen"
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryTopTen, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryTopTen,
		cacheKey,
		func() (interface{}, error) {
			topTen, err := extractors.ExtractTopTen(h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch top ten: %v", err))
			}
			return topTen, nil
		},
	)

	if err != nil {
		if fErr, ok := err.(*fiber.Error); ok {
			return fErr
		}
		return fiber.NewError(500, err.Error())
	}

	var result extractors.TopTenResult
	if err := json.Unmarshal(dataBytes, &result); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError,
			fmt.Sprintf("failed to parse top ten data: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"cached":  wasCached,
		"results": result,
	})
}
