package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// FilterHandler serves the filter endpoint.
type FilterHandler struct {
	baseHost     string
	cacheManager *cache.Manager
}

// NewFilterHandler creates a FilterHandler configured with the v1 provider host.
func NewFilterHandler(cacheManager *cache.Manager) *FilterHandler {
	return &FilterHandler{
		baseHost:     utils.GetV1BaseHost(),
		cacheManager: cacheManager,
	}
}

// Filter handles filter requests with various query parameters.
// @Summary      Filter anime
// @Description  Filter anime using the same query options as search (type, status, rating, genres, dates, etc.)
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        type      query     string  false  "Anime type (tv, movie, ova, etc.)"
// @Param        status    query     string  false  "Status (ongoing, completed, etc.)"
// @Param        rated     query     string  false  "Rating filter"
// @Param        score     query     string  false  "Minimum score"
// @Param        season    query     string  false  "Season filter"
// @Param        language  query     string  false  "Language (sub, dub)"
// @Param        genres    query     string  false  "Comma-separated genre IDs"  example(1,2,3)
// @Param        sort      query     string  false  "Sort order"
// @Param        sy        query     string  false  "Start year"
// @Param        sm        query     string  false  "Start month"
// @Param        sd        query     string  false  "Start day"
// @Param        ey        query     string  false  "End year"
// @Param        em        query     string  false  "End month"
// @Param        ed        query     string  false  "End day"
// @Param        keyword   query     string  false  "Search keyword"
// @Param        page      query     int     false  "Page number"  default(1)  example(1)
// @Success      200       {object}  SearchResponse  "Filter results"
// @Failure      404       {object}  ErrorResponse   "Requested page exceeds total available pages"
// @Failure      502       {object}  ErrorResponse   "Failed to filter"
// @Router       /filter [get]
func (h *FilterHandler) Filter(c *fiber.Ctx) error {
	params := make(map[string]string)

	// Extract all possible query parameters (only include non-empty values)
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
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		params["keyword"] = keyword
	}

	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}
	if page > 1 {
		params["page"] = strconv.Itoa(page)
	}

	// Generate cache key from params
	cacheKey := h.generateFilterCacheKey(params)
	ctx := c.Context()

	_, cacheErr := h.cacheManager.Get(ctx, cache.CategorySearch, cacheKey)
	wasCached := (cacheErr == nil)

	dataBytes, err := h.cacheManager.GetOrSet(
		ctx,
		cache.CategorySearch,
		cacheKey,
		func() (interface{}, error) {
			result, err := extractors.ExtractFilter(params, h.baseHost)
			if err != nil {
				if strings.Contains(err.Error(), "exceeds total available pages") {
					return nil, fiber.NewError(fiber.StatusNotFound, err.Error())
				}
				return nil, fiber.NewError(fiber.StatusBadGateway, fmt.Sprintf("failed to filter: %v", err))
			}

			if result.TotalPage > 0 && page > result.TotalPage {
				return nil, fiber.NewError(fiber.StatusNotFound, "requested page exceeds total available pages")
			}

			responseData := fiber.Map{
				"data":      result.Data,
				"totalPage": result.TotalPage,
				"page":      result.Page,
				"hasNext":   result.HasNext,
			}
			return responseData, nil
		},
	)

	if err != nil {
		// Handle Fiber errors
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

// generateFilterCacheKey generates a cache key from filter parameters
func (h *FilterHandler) generateFilterCacheKey(params map[string]string) string {
	key := "filter:"
	for k, v := range params {
		if v != "" {
			key += fmt.Sprintf("%s:%s:", k, v)
		}
	}
	// Hash the key to ensure it's not too long
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("filter:%s", hex.EncodeToString(hash[:]))
}
