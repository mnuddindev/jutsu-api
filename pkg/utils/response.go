package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
	})
}

// PaginatedResponse sends a paginated response
func PaginatedResponse(c *fiber.Ctx, statusCode int, message string, data interface{}, meta interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// OK sends a 200 OK response
func OK(c *fiber.Ctx, message string, data interface{}) error {
	return SuccessResponse(c, fiber.StatusOK, message, data)
}

// Created sends a 201 Created response
func Created(c *fiber.Ctx, message string, data interface{}) error {
	return SuccessResponse(c, fiber.StatusCreated, message, data)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message)
}

// NotFound sends a 404 Not Found response
func NotFound(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}

// UnprocessableEntity sends a 422 Unprocessable Entity response
func UnprocessableEntity(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnprocessableEntity, message)
}

