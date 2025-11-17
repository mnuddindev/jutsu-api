package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// QtipHandler serves the qtip endpoint.
type QtipHandler struct {
	baseHost string
}

// NewQtipHandler creates a QtipHandler configured with the v1 provider host.
func NewQtipHandler() *QtipHandler {
	return &QtipHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetQtip returns qtip data for an anime.
func (h *QtipHandler) GetQtip(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	qtip, err := extractors.ExtractQtip(id, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch qtip: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": qtip,
	})
}
