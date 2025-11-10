package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Services  map[string]ServiceInfo `json:"services"`
}

// ServiceInfo represents service health information
type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
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
	services := make(map[string]ServiceInfo)
	allHealthy := true

	// Check cache (optional - don't fail if cache is down)
	cacheStatus := "healthy"
	cacheMessage := ""
	if err := cache.HealthCheck(); err != nil {
		cacheStatus = "unhealthy"
		cacheMessage = err.Error()
		// Cache is optional, so we don't set allHealthy to false
	}
	services["cache"] = ServiceInfo{
		Status:  cacheStatus,
		Message: cacheMessage,
	}

	status := "healthy"
	statusCode := fiber.StatusOK
	if !allHealthy {
		status = "unhealthy"
		statusCode = fiber.StatusServiceUnavailable
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: utils.FormatTimeString(utils.Now()),
		Services:  services,
	}

	return c.Status(statusCode).JSON(response)
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

