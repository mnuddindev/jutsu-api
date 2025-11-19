package handler

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// SuggestionHandler serves the search suggestions endpoint.
type SuggestionHandler struct {
	baseHost string
}

// NewSuggestionHandler creates a SuggestionHandler configured with the v1 provider host.
func NewSuggestionHandler() *SuggestionHandler {
	return &SuggestionHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetSuggestions returns search suggestions for a keyword.
func (h *SuggestionHandler) GetSuggestions(c *fiber.Ctx) error {
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'keyword' is required")
	}

	// URL encode the keyword (same as Node.js)
	encodedKeyword := url.QueryEscape(keyword)

	// Generate cache key
	cacheKey := fmt.Sprintf("suggestion:%s", encodedKeyword)

	// Try to get from cache
	var cached []extractors.SuggestionItem
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && len(cached) > 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	suggestions, err := extractors.ExtractSuggestions(encodedKeyword, h.baseHost)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch suggestions: %v", err))
	}

	// Return empty array instead of nil
	if suggestions == nil {
		suggestions = []extractors.SuggestionItem{}
	}

	// Cache the response
	_ = helper.SetCachedData(cacheKey, suggestions, helper.SuggestionCacheTTL)

	return c.JSON(fiber.Map{
		"success": true,
		"results": suggestions,
	})
}
