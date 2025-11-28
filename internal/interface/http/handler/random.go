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

// GetRandomID returns only the random anime identifier.
// @Summary      Get random anime ID
// @Description  Returns a random anime identifier
// @Tags         Random
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Random anime ID"
// @Failure      502  {object}  ErrorResponse           "Bad Gateway"
// @Failure      503  {object}  ErrorResponse           "Service Unavailable"
// @Router       /random/id [get]
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
// @Summary      Get random anime
// @Description  Returns full information for a randomly selected anime
// @Tags         Random
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Random anime information"
// @Failure      502  {object}  ErrorResponse           "Bad Gateway"
// @Router       /random [get]
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
