package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// TopTenHandler serves the top ten anime endpoint.
type TopTenHandler struct {
	baseHost string
}

// NewTopTenHandler creates a TopTenHandler configured with the v1 provider host.
func NewTopTenHandler() *TopTenHandler {
	return &TopTenHandler{
		baseHost: utils.GetV1BaseHost(),
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
	var cached extractors.TopTenResult
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached != nil && len(cached) > 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	topTen, err := extractors.ExtractTopTen(h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch top ten: %v", err))
	}

	_ = helper.SetCachedData(cacheKey, topTen, helper.TopTenCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": topTen,
	})
}
