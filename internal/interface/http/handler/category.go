package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// CategoryHandler serves category/genre listing endpoints.
type CategoryHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewCategoryHandler creates a CategoryHandler configured with the v1 provider host.
func NewCategoryHandler(cacheManager *cache.Manager) *CategoryHandler {
	return &CategoryHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// GetCategory returns paginated anime rows for a category/genre route.
// @Summary      Get category or genre listing
// @Description  Returns paginated anime rows for a given category or genre route (e.g. /api/genre/action, /api/top-airing)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        page   query     int   false  "Page number (must be >= 1)"  default(1)  example(1)
// @Success      200    {object}  CategoryResponse  "Category results"
// @Failure      400    {object}  ErrorResponse     "Invalid page parameter"
// @Failure      404    {object}  ErrorResponse     "Requested page exceeds total available pages"
// @Failure      502    {object}  ErrorResponse     "Failed to load category data"
// @Router       /genre/{slug} [get]
// @Router       /top-airing [get]
// @Router       /most-popular [get]
// @Router       /most-favorite [get]
// @Router       /completed [get]
// @Router       /recently-updated [get]
// @Router       /recently-added [get]
// @Router       /top-upcoming [get]
// @Router       /subbed-anime [get]
// @Router       /dubbed-anime [get]
// @Router       /movie [get]
// @Router       /special [get]
// @Router       /ova [get]
// @Router       /ona [get]
// @Router       /tv [get]
// @Router       /az-list [get]
// @Router       /az-list/{letter} [get]
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	// Extract route type from path (e.g., "/api/genre/action" -> "genre/action")
	path := c.Path()
	routeType := strings.TrimPrefix(path, "/api/")
	// Handle martial-arts typo fix
	if routeType == "genre/martial-arts" {
		routeType = "genre/marial-arts"
	}

	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'page' must be a positive integer")
	}

	cacheKey := fmt.Sprintf("%s_page_%d", strings.ReplaceAll(routeType, "/", "_"), page)
	ctx := c.Context()

	// Check if cached (for response field)
	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryGenre, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryGenre,
		cacheKey,
		func() (interface{}, error) {
			result, err := extractors.ExtractCategory(routeType, page, h.baseHost)
			if err != nil {
				if strings.Contains(err.Error(), "exceeds total available pages") {
					return nil, fiber.NewError(fiber.StatusNotFound, err.Error())
				}
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to load category data: %v", err))
			}

			if result.TotalPages > 0 && page > result.TotalPages {
				return nil, fiber.NewError(fiber.StatusNotFound, "requested page exceeds total available pages")
			}
			return result, nil
		},
	)

	if err != nil {
		// Handle Fiber errors
		if fErr, ok := err.(*fiber.Error); ok {
			return fErr
		}
		return fiber.NewError(500, err.Error())
	}

	var result extractors.CategoryResult
	json.Unmarshal(dataBytes, &result)

	return c.JSON(fiber.Map{
		"success": true,
		"cached":  wasCached,
		"results": result,
	})
}
