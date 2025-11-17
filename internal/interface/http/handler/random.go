package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// RandomHandler serves random anime endpoints.
type RandomHandler struct {
	baseHost string
}

// NewRandomHandler builds a handler configured with the v1 provider host.
func NewRandomHandler() *RandomHandler {
	return &RandomHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetRandomID returns only the random anime identifier (matching the Node API behavior).
func (h *RandomHandler) GetRandomID(c *fiber.Ctx) error {
	id, err := extractors.ExtractRandomID(h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
	if strings.TrimSpace(id) == "" {
		return fiber.NewError(fiber.StatusServiceUnavailable, "no random id could be determined")
	}
	return c.JSON(fiber.Map{
		"success": true,
		"results": fiber.Map{"id": id},
	})
}

// GetRandom returns the full anime information for a randomly selected title.
func (h *RandomHandler) GetRandom(c *fiber.Ctx) error {
	info, err := extractors.ExtractRandom(h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
	return c.JSON(fiber.Map{
		"success": true,
		"results": info,
	})
}
