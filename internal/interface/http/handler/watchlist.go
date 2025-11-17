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

	// Restructure response to match Node.js format
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

	return c.JSON(fiber.Map{
		"success": true,
		"results": fiber.Map{
			"totalPages": result.TotalPages,
			"data":       watchlistData,
		},
	})
}
