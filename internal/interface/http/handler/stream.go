package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// StreamHandler serves the stream and fallback endpoints.
type StreamHandler struct {
	baseHost string
}

// NewStreamHandler creates a StreamHandler that targets the v1 provider host.
func NewStreamHandler() *StreamHandler {
	return &StreamHandler{
		baseHost: utils.GetV1BaseHost(),
	}
}

// GetStream resolves the default streaming info for an episode/server pair.
// @Summary      Get streaming information
// @Description  Returns streaming information for an episode with specified server and type
// @Tags         Streaming
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Anime ID"  example(frieren-beyond-journeys-end-18542)
// @Param        ep       query     string  true  "Episode ID"  example(107257)
// @Param        server   query     string  false  "Server name (e.g., vidcloud, hd-1)"  example(hd-1)
// @Param        type     query     string  false  "Stream type (sub or dub)"  example(sub)
// @Success      200      {object}  map[string]interface{}  "Streaming information"
// @Failure      400      {object}  map[string]interface{}  "Bad Request"
// @Failure      502      {object}  map[string]interface{}  "Bad Gateway"
// @Router       /stream/{id} [get]
func (h *StreamHandler) GetStream(c *fiber.Ctx) error {
	return h.handleStreamRequest(c, false)
}

// GetStreamFallback resolves the fallback streaming info for an episode/server pair.
// @Summary      Get fallback streaming information
// @Description  Returns fallback streaming information when primary stream fails
// @Tags         Streaming
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Anime ID"  example(frieren-beyond-journeys-end-18542)
// @Param        ep       query     string  true  "Episode ID"  example(107257)
// @Param        server   query     string  false  "Server name"  example(hd-1)
// @Param        type     query     string  false  "Stream type (sub or dub)"  example(sub)
// @Success      200      {object}  map[string]interface{}  "Fallback streaming information"
// @Failure      400      {object}  map[string]interface{}  "Bad Request"
// @Failure      503      {object}  map[string]interface{}  "Service Unavailable"
// @Router       /stream/fallback/{id} [get]
func (h *StreamHandler) GetStreamFallback(c *fiber.Ctx) error {
	return h.handleStreamRequest(c, true)
}

func (h *StreamHandler) handleStreamRequest(c *fiber.Ctx, fallback bool) error {
	// Get anime ID from path parameter
	animeID := strings.TrimSpace(c.Params("id"))
	if animeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "path parameter 'id' is required")
	}

	// Get episode ID from query parameter
	episodeID := strings.TrimSpace(c.Query("ep"))
	if episodeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query parameter 'ep' is required")
	}

	serverName := c.Query("server")
	streamType := c.Query("type")

	// Generate cache key based on all parameters
	cacheKey := h.generateStreamCacheKey(animeID, episodeID, serverName, streamType, fallback)

	// Try to get from cache
	var cached map[string]interface{}
	if err := helper.GetCachedData(cacheKey, &cached); err == nil && cached != nil && len(cached) > 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"results": cached,
			"cached":  true,
		})
	}

	// Build full ID format: anime-id?ep=episode-id (matching Node.js format)
	// The extractor will extract just the episode ID for servers, but use full format for decryption
	fullID := fmt.Sprintf("%s?ep=%s", animeID, episodeID)

	streamInfo, err := extractors.ExtractStreamingInfo(fullID, serverName, streamType, fallback, h.baseHost)
	if err != nil {
		status := fiber.StatusBadGateway
		if fallback {
			status = fiber.StatusServiceUnavailable
		}
		return fiber.NewError(status, fmt.Sprintf("failed to resolve streaming info: %v", err))
	}

	// Transform response to match expected format
	response := h.transformStreamResponse(streamInfo)

	// Cache the response (only cache successful responses with data)
	streamingLinks, _ := response["streamingLink"].([]interface{})
	servers, _ := response["servers"].([]map[string]interface{})
	if len(streamingLinks) > 0 || len(servers) > 0 {
		_ = helper.SetCachedData(cacheKey, response, helper.StreamInfoCacheTTL)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"results": response,
	})
}

// generateStreamCacheKey generates a unique cache key for stream requests
func (h *StreamHandler) generateStreamCacheKey(animeID, episodeID, serverName, streamType string, fallback bool) string {
	key := fmt.Sprintf("stream:%s:%s", animeID, episodeID)
	if serverName != "" {
		key += ":" + serverName
	}
	if streamType != "" {
		key += ":" + streamType
	}
	if fallback {
		key += ":fallback"
	}
	// Hash the key to ensure it's not too long
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("stream:%s", hex.EncodeToString(hash[:]))
}

// transformStreamResponse transforms the extractor response to match the expected API format
func (h *StreamHandler) transformStreamResponse(streamInfo extractors.StreamingInfo) map[string]interface{} {
	// Transform streamingLink to array format
	streamingLinkArray := []interface{}{}
	if streamInfo.StreamingLink.ID != "" || streamInfo.StreamingLink.Link.File != "" {
		// Convert ID to number if possible
		var idNum interface{} = streamInfo.StreamingLink.ID
		if idNumStr, err := parseNumber(streamInfo.StreamingLink.ID); err == nil {
			idNum = idNumStr
		}

		streamingLinkItem := map[string]interface{}{
			"id":     idNum,
			"type":   streamInfo.StreamingLink.Type,
			"link":   streamInfo.StreamingLink.Link,
			"tracks": streamInfo.StreamingLink.Tracks,
			"intro":  streamInfo.StreamingLink.Intro,
			"outro":  streamInfo.StreamingLink.Outro,
			"server": streamInfo.StreamingLink.Server,
		}
		streamingLinkArray = append(streamingLinkArray, streamingLinkItem)
	}

	// Transform servers to match expected format
	serversArray := []map[string]interface{}{}
	for _, server := range streamInfo.Servers {
		// Convert data_id, server_id to numbers if possible
		var dataIDNum interface{} = server.DataID
		var serverIDNum interface{} = server.ServerID
		if num, err := parseNumber(server.DataID); err == nil {
			dataIDNum = num
		}
		if num, err := parseNumber(server.ServerID); err == nil {
			serverIDNum = num
		}

		serversArray = append(serversArray, map[string]interface{}{
			"type":        server.Type,
			"data_id":     dataIDNum,
			"server_id":   serverIDNum,
			"server_name": server.ServerName,
		})
	}

	return map[string]interface{}{
		"streamingLink": streamingLinkArray,
		"servers":       serversArray,
	}
}

// parseNumber attempts to parse a string as a number (int or float)
func parseNumber(s string) (interface{}, error) {
	// Try parsing as int first
	if i, err := strconv.Atoi(s); err == nil {
		return i, nil
	}
	// Try parsing as float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	return nil, fmt.Errorf("not a number")
}
