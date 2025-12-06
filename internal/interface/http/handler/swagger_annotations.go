package handler

// This file contains Swagger annotation helpers and common response types
// Used for generating comprehensive Swagger documentation

// Common response types for Swagger documentation

// SuccessResponse represents a successful API response
// @Description Standard success response
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Results interface{} `json:"results"`
	Cached  bool        `json:"cached,omitempty" example:"false"`
}

// ErrorResponse represents an error API response
// @Description Standard error response
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Error message here"`
}

// HomeResponse represents home page data response
type HomeResponse struct {
	Success bool                   `json:"success" example:"true"`
	Results map[string]interface{} `json:"results"`
}

// AnimeInfoResponse represents anime information response
type AnimeInfoResponse struct {
	Success bool                   `json:"success" example:"true"`
	Results map[string]interface{} `json:"results"`
	Data    interface{}            `json:"data"`
	Seasons interface{}            `json:"seasons"`
}

// StreamingInfoResponse represents streaming information response
type StreamingInfoResponse struct {
	Success bool                   `json:"success" example:"true"`
	Results map[string]interface{} `json:"results"`
}

// CategoryResponse represents category listing response
type CategoryResponse struct {
	Success bool `json:"success" example:"true"`
	Results struct {
		TotalPages int         `json:"totalPages" example:"10"`
		Data       interface{} `json:"data"`
	} `json:"results"`
}

// SearchResponse represents search results response
type SearchResponse struct {
	Success bool `json:"success" example:"true"`
	Results struct {
		Data      interface{} `json:"data"`
		TotalPage int         `json:"totalPage" example:"5"`
	} `json:"results"`
}
