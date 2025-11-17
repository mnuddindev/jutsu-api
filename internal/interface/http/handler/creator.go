package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// CreatorHandler serves producer/studio listing endpoints.
type CreatorHandler struct {
	baseHost string
}

// NewCreatorHandler builds a handler configured with the v1 provider host.
func NewCreatorHandler() *CreatorHandler {
	return &CreatorHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetProducer returns paginated anime rows for a producer slug.
func (h *CreatorHandler) GetProducer(c *fiber.Ctx) error {
	return h.handleCreatorRequest(c, "producer")
}

// GetStudio returns paginated anime rows for a studio slug.
func (h *CreatorHandler) GetStudio(c *fiber.Ctx) error {
	return h.handleCreatorRequest(c, "studio")
}

func (h *CreatorHandler) handleCreatorRequest(c *fiber.Ctx, prefix string) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'page' must be a positive integer")
	}

	cacheKey := fmt.Sprintf("creator:%s:%s:%d", prefix, id, page)
	var cached extractors.CategoryResult
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached.Data != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	var result extractors.CategoryResult
	switch prefix {
	case "producer":
		result, err = extractors.ExtractProducer(id, page, h.baseHost)
	default:
		result, err = extractors.ExtractStudio(id, page, h.baseHost)
	}
	if err != nil {
		if err == extractors.ErrCreatorPageOutOfRange {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		if err == extractors.ErrCreatorIDRequired {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to load %s data: %v", prefix, err))
	}

	_ = helper.SetCachedData(cacheKey, result, helper.CreatorCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": result,
	})
}
