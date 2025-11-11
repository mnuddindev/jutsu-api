package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/config"
	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// CacheCheck represents cache health
type CacheCheck struct {
	Status         string `json:"status"`
	ResponseTimeMs int64  `json:"response_time_ms"`
}

// ExternalSourcesCheck represents external source health
type ExternalSourcesCheck struct {
	AnimeProvider map[string]string `json:"anime_provider"`
}

// Checks groups all checks
type Checks struct {
	Cache           CacheCheck           `json:"cache"`
	ExternalSources ExternalSourcesCheck `json:"external_sources"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Uptime    string `json:"uptime"`
	Checks    Checks `json:"checks"`
	Timestamp string `json:"timestamp"`
}

// Health checks the health of the application and its services
// @Summary Health check
// @Description Check the health of the application and its services
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	cfg := config.Cfg
	if cfg == nil {
		// Fallback if not initialized
		if loaded, err := config.LoadConfig(); err == nil {
			cfg = loaded
		}
	}

	// Cache check
	cacheStart := time.Now()
	cacheStatus := "ok"
	if err := cache.HealthCheck(); err != nil {
		cacheStatus = "unhealthy"
	}
	cacheCheck := CacheCheck{
		Status:         cacheStatus,
		ResponseTimeMs: time.Since(cacheStart).Milliseconds(),
	}

	// External providers check with Redis cache (24h)
	providerStatus := map[string]string{}
	cached := false
	if cache.Client != nil {
		if err := cache.Get("providers:status", &providerStatus); err == nil && len(providerStatus) > 0 {
			cached = true
		}
	}
	if !cached {
		providers := utils.GetBaseProviders()
		for _, p := range providers {
			providerStatus[p] = probeURL(p)
		}
		// cache for 24h
		_ = cache.Set("providers:status", providerStatus, 24*time.Hour)
	}

	external := ExternalSourcesCheck{AnimeProvider: providerStatus}

	resp := HealthResponse{
		Status:    "ok",
		Service:   cfg.App.Name,
		Version:   cfg.App.Version,
		Uptime:    utils.GetUptimeString(),
		Checks:    Checks{Cache: cacheCheck, ExternalSources: external},
		Timestamp: utils.FormatTimeString(utils.Now()),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Ready checks if the application is ready to serve traffic
// @Summary Readiness check
// @Description Check if the application is ready to serve traffic
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ready [get]
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"status": "ready",
	})
}

// Live checks if the application is alive
// @Summary Liveness check
// @Description Check if the application is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /live [get]
func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"status": "alive",
	})
}

// probeURL checks if a provider is reachable and returns "live" or "down"
func probeURL(hostOrURL string) string {
	url := hostOrURL
	if !strings.HasPrefix(strings.ToLower(url), "http://") && !strings.HasPrefix(strings.ToLower(url), "https://") {
		url = "https://" + hostOrURL
	}
	client := &http.Client{Timeout: 2 * time.Second}
	// Try HEAD first
	req, _ := http.NewRequest(http.MethodHead, url, nil)
	resp, err := client.Do(req)
	if err == nil && resp != nil && resp.StatusCode < 500 {
		return "live"
	}
	// Fallback to GET
	resp, err = client.Get(url)
	if err == nil && resp != nil && resp.StatusCode < 500 {
		return "live"
	}
	return "down"
}
