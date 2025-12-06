package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// ScheduleHandler serves the schedule endpoint.
type ScheduleHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewScheduleHandler creates a ScheduleHandler configured with the v1 provider host.
func NewScheduleHandler(cacheManager *cache.Manager) *ScheduleHandler {
	return &ScheduleHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
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
		date = time.Now().Format("2006-01-02")
	}

	tzOffsetParam := c.Query("tzOffset", "-330")
	tzOffset, err := strconv.Atoi(tzOffsetParam)
	if err != nil {
		tzOffset = -330
	}

	cacheKey := fmt.Sprintf("schedule:%s:%d", date, tzOffset)
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategorySchedule, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategorySchedule,
		cacheKey,
		func() (interface{}, error) {
			schedule, err := extractors.ExtractSchedule(date, tzOffset, h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch schedule: %v", err))
			}
			return schedule, nil
		},
	)

	if err != nil {
		if fErr, ok := err.(*fiber.Error); ok {
			return fErr
		}
		return fiber.NewError(500, err.Error())
	}

	var result map[string]interface{}
	json.Unmarshal(dataBytes, &result)

	return c.JSON(fiber.Map{
		"success": true,
		"cached":  wasCached,
		"results": result,
	})
}
