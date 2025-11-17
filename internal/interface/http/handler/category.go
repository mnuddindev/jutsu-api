package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

const categoryCacheTTL = 30 * time.Minute

// CategoryHandler serves category/genre listing endpoints.
type CategoryHandler struct {
	baseHost string
}

// NewCategoryHandler creates a CategoryHandler configured with the v1 provider host.
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetCategory returns paginated anime rows for a category/genre route.
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	// Extract route type from path (e.g., "/api/genre/action" -> "genre/action")
	path := c.Path()
	routeType := strings.TrimPrefix(path, "/api/")
	// Handle martial-arts typo fix (same as Node.js)
	if routeType == "genre/martial-arts" {
		routeType = "genre/marial-arts"
	}

	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'page' must be a positive integer")
	}

	cacheKey := fmt.Sprintf("%s_page_%d", strings.ReplaceAll(routeType, "/", "_"), page)
	var cached extractors.CategoryResult
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached.Data != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	result, err := extractors.ExtractCategory(routeType, page, h.baseHost)
	if err != nil {
		if strings.Contains(err.Error(), "exceeds total available pages") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to load category data: %v", err))
	}

	if result.TotalPages > 0 && page > result.TotalPages {
		return fiber.NewError(fiber.StatusNotFound, "requested page exceeds total available pages")
	}

	_ = helper.SetCachedData(cacheKey, result, categoryCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": result,
	})
}
