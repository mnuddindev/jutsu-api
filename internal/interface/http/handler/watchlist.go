package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// WatchlistHandler serves the watchlist endpoint.
type WatchlistHandler struct {
	baseHost string
}

// NewWatchlistHandler creates a WatchlistHandler configured with the v1 provider host.
func NewWatchlistHandler() *WatchlistHandler {
	return &WatchlistHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetWatchlist returns a user's watchlist with pagination.
// @Summary      Get user watchlist
// @Description  Returns a paginated watchlist for a specific user
// @Tags         Watchlist
// @Accept       json
// @Produce      json
// @Param        userId   path      string  true  "User identifier"  example(user123)
// @Param        page     path      int     false "Page number (must be >= 1)"  default(1)  example(1)
// @Success      200      {object}  SuccessResponse  "Watchlist data"
// @Failure      400      {object}  ErrorResponse    "Missing or invalid parameters"
// @Failure      502      {object}  ErrorResponse    "Failed to fetch watchlist"
// @Router       /watchlist/{userId} [get]
// @Router       /watchlist/{userId}/{page} [get]
func (h *WatchlistHandler) GetWatchlist(c *fiber.Ctx) error {
	userID := strings.TrimSpace(c.Params("userId"))
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'userId' is required")
	}

	pageParam := c.Params("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	result, err := extractors.ExtractWatchlist(userID, page, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch watchlist: %v", err))
	}

	watchlistData := make([]map[string]interface{}, 0, len(result.Watchlist))
	for _, item := range result.Watchlist {
		watchlistData = append(watchlistData, map[string]interface{}{
			"id":       item.ID,
			"poster":   item.Poster,
			"title":    item.Title,
			"duration": item.Duration,
			"type":     item.Type,
			"subCount": item.SubCount,
			"dubCount": item.DubCount,
			"link":     item.Link,
			"showType": item.ShowType,
			"tvInfo": map[string]interface{}{
				"showType": item.TVInfo.ShowType,
				"duration": item.TVInfo.Duration,
				"sub":      item.TVInfo.Sub,
				"dub":      item.TVInfo.Dub,
			},
		})
	}

	responseData := fiber.Map{
		"totalPages": result.TotalPages,
		"data":       watchlistData,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": responseData,
	})
}
