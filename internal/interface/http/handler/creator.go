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

// CreatorHandler serves producer/studio listing endpoints.
type CreatorHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewCreatorHandler builds a handler configured with the v1 provider host.
func NewCreatorHandler(cacheManager *cache.Manager) *CreatorHandler {
	return &CreatorHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// GetProducer returns paginated anime rows for a producer slug.
// @Summary      Get anime of specific producers
// @Description  Returns paginated anime rows for a producer identifier
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Producer ID or slug"  example(1)
// @Param        page   query     int     false "Page number (must be >= 1)"  default(1)  example(1)
// @Success      200    {object}  CategoryResponse  "Producer anime list"
// @Failure      400    {object}  ErrorResponse     "Invalid parameters"
// @Failure      404    {object}  ErrorResponse     "Page out of range"
// @Failure      502    {object}  ErrorResponse     "Failed to load producer data"
// @Router       /producer/{id} [get]
func (h *CreatorHandler) GetProducer(c *fiber.Ctx) error {
	return h.handleCreatorRequest(c, "producer")
}

// GetStudio returns paginated anime rows for a studio slug.
// @Summary      Get anime of specific studios
// @Description  Returns paginated anime rows for a studio identifier
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Studio ID or slug"  example(1)
// @Param        page   query     int     false "Page number (must be >= 1)"  default(1)  example(1)
// @Success      200    {object}  CategoryResponse  "Studio anime list"
// @Failure      400    {object}  ErrorResponse     "Invalid parameters"
// @Failure      404    {object}  ErrorResponse     "Page out of range"
// @Failure      502    {object}  ErrorResponse     "Failed to load studio data"
// @Router       /studio/{id} [get]
func (h *CreatorHandler) GetStudio(c *fiber.Ctx) error {
	return h.handleCreatorRequest(c, "studio")
}

func (h *CreatorHandler) handleCreatorRequest(c *fiber.Ctx, prefix string) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	pageParam := c.Query("page", "1")
	ctx := c.Context()

	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'page' must be a positive integer")
	}

	cacheKey := fmt.Sprintf("creator:%s:%s:%d", prefix, id, page)

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategoryProducer, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategoryProducer,
		cacheKey,
		func() (interface{}, error) {
			var result extractors.CategoryResult
			switch prefix {
			case "producer":
				result, err = extractors.ExtractProducer(id, page, h.baseHost)
			default:
				result, err = extractors.ExtractStudio(id, page, h.baseHost)
			}
			if err != nil {
				if err == extractors.ErrCreatorPageOutOfRange {
					return nil, fiber.NewError(fiber.StatusNotFound, err.Error())
				}
				if err == extractors.ErrCreatorIDRequired {
					return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
				}
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to load %s data: %v", prefix, err))
			}
			return result, nil
		},
	)

	if err != nil {
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
