package validation

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"

	appLogger "github.com/mnuddindev/jutsu-api/internal/infrastructure/logger"
)

var (
	Validate   *validator.Validate
	Translator ut.Translator
)

// InitValidator initializes the validator
func InitValidator() error {
	// Create validator instance
	Validate = validator.New()

	// Register custom tag name function to use json tags
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Setup English translator
	en := en.New()
	uni := ut.New(en, en)
	Translator, _ = uni.GetTranslator("en")

	// Register default translations
	if err := enTranslations.RegisterDefaultTranslations(Validate, Translator); err != nil {
		return fmt.Errorf("failed to register default translations: %w", err)
	}

	// Register custom validators
	registerCustomValidators()

	return nil
}

// ValidateStruct validates a struct and returns errors if any
func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError

	err := Validate.Struct(s)
	if err != nil {
		validatorErrors := err.(validator.ValidationErrors)
		for _, e := range validatorErrors {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Tag:     e.Tag(),
				Value:   fmt.Sprintf("%v", e.Value()),
				Message: e.Translate(Translator),
			})
		}
	}

	return errors
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// ValidateRequest validates a request and returns a fiber error if validation fails
func ValidateRequest(c *fiber.Ctx, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		appLogger.Error("Failed to bind request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	errors := ValidateStruct(req)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  errors,
		})
	}

	return nil
}

// registerCustomValidators registers custom validation functions
func registerCustomValidators() {
	// Add custom validators here
	// Example:
	// Validate.RegisterValidation("custom_tag", customValidatorFunc)
}

