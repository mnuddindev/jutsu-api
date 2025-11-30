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

	// Try to get from cache first
	var cached extractors.TopTenResult
	if err := helper.GetCachedData(cacheKey, &cached); err == nil {
		// Check if we have valid cached data
		hasData := false
		if cached != nil {
			for _, items := range cached {
				if len(items) > 0 {
					hasData = true
					break
				}
			}
		}

		if hasData {
			return c.JSON(fiber.Map{
				"success": true,
				"results": cached,
				"cached":  true,
			})
		}
	}

	// If no cache or invalid cache, fetch fresh data
	topTen, err := extractors.ExtractTopTen(h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch top ten: %v", err))
	}

	// Cache the fresh data
	if err := helper.SetCachedData(cacheKey, topTen, helper.TopTenCacheTTL); err != nil {
		// Log but don't fail the request
		fmt.Printf("Failed to cache top ten data: %v\n", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": topTen,
	})
}
