package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// SearchHandler serves the search endpoint.
type SearchHandler struct {
	baseHost string
}

// NewSearchHandler creates a SearchHandler configured with the v1 provider host.
func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// Search handles search requests with various query parameters.
func (h *SearchHandler) Search(c *fiber.Ctx) error {
	params := make(map[string]string)

	// Extract all possible query parameters
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		params["keyword"] = keyword
	}
	if typ := strings.TrimSpace(c.Query("type")); typ != "" {
		params["type"] = typ
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		params["status"] = status
	}
	if rated := strings.TrimSpace(c.Query("rated")); rated != "" {
		params["rated"] = rated
	}
	if score := strings.TrimSpace(c.Query("score")); score != "" {
		params["score"] = score
	}
	if season := strings.TrimSpace(c.Query("season")); season != "" {
		params["season"] = season
	}
	if language := strings.TrimSpace(c.Query("language")); language != "" {
		params["language"] = language
	}
	if genres := strings.TrimSpace(c.Query("genres")); genres != "" {
		params["genres"] = genres
	}
	if sort := strings.TrimSpace(c.Query("sort")); sort != "" {
		params["sort"] = sort
	}
	if sy := strings.TrimSpace(c.Query("sy")); sy != "" {
		params["sy"] = sy
	}
	if sm := strings.TrimSpace(c.Query("sm")); sm != "" {
		params["sm"] = sm
	}
	if sd := strings.TrimSpace(c.Query("sd")); sd != "" {
		params["sd"] = sd
	}
	if ey := strings.TrimSpace(c.Query("ey")); ey != "" {
		params["ey"] = ey
	}
	if em := strings.TrimSpace(c.Query("em")); em != "" {
		params["em"] = em
	}
	if ed := strings.TrimSpace(c.Query("ed")); ed != "" {
		params["ed"] = ed
	}

	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}
	if page > 1 {
		params["page"] = strconv.Itoa(page)
	}

	result, err := extractors.ExtractSearch(params, h.baseHost)
	if err != nil {
		if strings.Contains(err.Error(), "exceeds total available pages") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to search: %v", err))
	}

	if result.TotalPage > 0 && page > result.TotalPage {
		return fiber.NewError(fiber.StatusNotFound, "requested page exceeds total available pages")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"results": fiber.Map{
			"data":      result.Data,
			"totalPage": result.TotalPage,
		},
	})
}
