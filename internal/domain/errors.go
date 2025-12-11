package domain

import (
	"errors"
	"fmt"
)

// Authentication errors
var (
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserNotActive         = errors.New("user account is not active")
	ErrEmailNotVerified      = errors.New("email not verified")
)

// Token errors
var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrTokenRevoked       = errors.New("token has been revoked")
	ErrTokenNotFound      = errors.New("token not found")
	ErrWrongTokenType     = errors.New("wrong token type")
	ErrTokenBlacklisted   = errors.New("token is blacklisted")
	ErrInvalidTokenFormat = errors.New("invalid token format")
)

// Validation errors
var (
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrWeakPassword     = errors.New("password is too weak")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong  = errors.New("password must be less than 72 characters")
	ErrInvalidUsername  = errors.New("invalid username format")
	ErrInvalidInput     = errors.New("invalid input")
	ErrValidationFailed = errors.New("validation failed")
)

// Permission errors
var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden: insufficient permissions")
	ErrInvalidRole      = errors.New("invalid role")
	ErrPermissionDenied = errors.New("permission denied")
)

// Rate limiting errors
var (
	ErrTooManyRequests   = errors.New("too many requests")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// Database errors
var (
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")
	ErrRecordNotFound     = errors.New("record not found")
	ErrDuplicateKey       = errors.New("duplicate key violation")
)

// Cache errors
var (
	ErrCacheConnection = errors.New("cache connection error")
	ErrCacheMiss       = errors.New("cache miss")
	ErrCacheSet        = errors.New("cache set error")
)

// AppError represents an application error with additional context
type AppError struct {
	Err        error
	Message    string
	Code       string
	StatusCode int
	Details    map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(err error, message string, code string, statusCode int) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
		Details:    make(map[string]interface{}),
	}
}

// WithDetail adds a detail to the error
func (e *AppError) WithDetail(key string, value interface{}) *AppError {
	e.Details[key] = value
	return e
}

// ValidationError represents a validation error with field-specific messages
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation error"
	}
	return fmt.Sprintf("validation failed: %s", ve[0].Message)
}

// HasField checks if a specific field has an error
func (ve ValidationErrors) HasField(field string) bool {
	for _, err := range ve {
		if err.Field == field {
			return true
		}
	}
	return false
}

// GetFieldError returns the error message for a specific field
func (ve ValidationErrors) GetFieldError(field string) string {
	for _, err := range ve {
		if err.Field == field {
			return err.Message
		}
	}
	return ""
}

// ToMap converts validation errors to a map
func (ve ValidationErrors) ToMap() map[string]string {
	result := make(map[string]string)
	for _, err := range ve {
		result[err.Field] = err.Message
	}
	return result
}
