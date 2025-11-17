package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// ScheduleHandler serves the schedule endpoint.
type ScheduleHandler struct {
	baseHost string
}

// NewScheduleHandler creates a ScheduleHandler configured with the v1 provider host.
func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetSchedule returns the schedule for a specific date.
// @Summary      Get anime schedule
// @Description  Returns anime schedule for a specific date
// @Tags         Schedule
// @Accept       json
// @Produce      json
// @Param        date      query     string  false  "Date in YYYY-MM-DD format (default: today)"  example(2025-01-18)
// @Param        tzOffset  query     int     false  "Timezone offset in minutes (default: -330)"  default(-330)  example(-330)
// @Success      200       {object}  map[string]interface{}  "Schedule data"
// @Failure      502       {object}  ErrorResponse           "Bad Gateway"
// @Router       /schedule [get]
func (h *ScheduleHandler) GetSchedule(c *fiber.Ctx) error {
	date := strings.TrimSpace(c.Query("date"))
	if date == "" {
		// Default to today's date (same as Node.js)
		date = time.Now().Format("2006-01-02")
	}

	tzOffsetParam := c.Query("tzOffset", "-330")
	tzOffset, err := strconv.Atoi(tzOffsetParam)
	if err != nil {
		tzOffset = -330 // Default timezone offset
	}

	schedule, err := extractors.ExtractSchedule(date, tzOffset, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch schedule: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": schedule,
	})
}
