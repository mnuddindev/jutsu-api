package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// TopSearchHandler serves the top search endpoint.
type TopSearchHandler struct {
	baseHost string
}

// NewTopSearchHandler creates a TopSearchHandler configured with the v1 provider host.
func NewTopSearchHandler() *TopSearchHandler {
	return &TopSearchHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetTopSearch returns top search keywords.
func (h *TopSearchHandler) GetTopSearch(c *fiber.Ctx) error {
	topSearch, err := extractors.ExtractTopSearch(h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch top search: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": topSearch,
	})
}
