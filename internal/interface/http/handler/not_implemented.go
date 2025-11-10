package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// NotImplementedHandler handles routes that are not yet implemented
type NotImplementedHandler struct{}

// NewNotImplementedHandler creates a new not implemented handler
func NewNotImplementedHandler() *NotImplementedHandler {
	return &NotImplementedHandler{}
}

// NotImplementedResponse represents the not implemented response
type NotImplementedResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Health  map[string]interface{} `json:"health,omitempty"`
}

// NotImplemented returns a not implemented response
func (h *NotImplementedHandler) NotImplemented(c *fiber.Ctx) error {
	// Optionally include health status
	includeHealth := c.Query("health", "false") == "true"
	
	response := NotImplementedResponse{
		Success: false,
		Message: "This endpoint is not yet implemented. It will be available soon.",
	}

	if includeHealth {
		health := make(map[string]interface{})
		
		// Check cache status
		cacheStatus := "healthy"
		if err := cache.HealthCheck(); err != nil {
			cacheStatus = "unhealthy"
		}
		health["cache"] = map[string]string{
			"status": cacheStatus,
		}
		
		health["timestamp"] = utils.FormatTimeString(utils.Now())
		response.Health = health
	}

	return c.Status(fiber.StatusNotImplemented).JSON(response)
}

