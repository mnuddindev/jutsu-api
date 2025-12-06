package handler

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// SuggestionHandler serves the search suggestions endpoint.
type SuggestionHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewSuggestionHandler creates a SuggestionHandler configured with the v1 provider host.
func NewSuggestionHandler(cacheManager *cache.Manager) *SuggestionHandler {
	return &SuggestionHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// GetSuggestions returns search suggestions for a keyword.
// @Summary      Get search suggestions
// @Description  Returns search suggestions for a given keyword
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword  query     string  true  "Search keyword"  example(nar)
// @Success      200      {object}  SuccessResponse  "Search suggestions"
// @Failure      400      {object}  ErrorResponse    "Missing keyword"
// @Failure      502      {object}  ErrorResponse    "Failed to fetch suggestions"
// @Router       /search/suggest [get]
func (h *SuggestionHandler) GetSuggestions(c *fiber.Ctx) error {
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'keyword' is required")
	}

	encodedKeyword := url.QueryEscape(keyword)

	cacheKey := fmt.Sprintf("suggestion:%s", encodedKeyword)
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategorySuggest, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategorySuggest,
		cacheKey,
		func() (interface{}, error) {
			suggestions, err := extractors.ExtractSuggestions(encodedKeyword, h.baseHost)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to fetch suggestions: %v", err))
			}

			if suggestions == nil {
				suggestions = []extractors.SuggestionItem{}
			}
			return suggestions, nil
		},
	)

	if err != nil {
		if fErr, ok := err.(*fiber.Error); ok {
			return fErr
		}
		return fiber.NewError(500, err.Error())
	}

	var result []extractors.SuggestionItem
	json.Unmarshal(dataBytes, &result)

	return c.JSON(fiber.Map{
		"success": true,
		"cached":  wasCached,
		"results": result,
	})
}
