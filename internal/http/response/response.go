package response

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/logger"
)

// Response wraps all API responses
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
}

type ErrorData struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Success sends successful response
func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends error response
func Error(c *fiber.Ctx, statusCode int, code, message string, details map[string]interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ValidationError sends validation error with field-level errors
func ValidationError(c *fiber.Ctx, fieldErrors map[string]string) error {
	details := make(map[string]interface{})
	details["fields"] = fieldErrors

	return Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
}

// HandleError processes AppError and sends appropriate response
func HandleError(c *fiber.Ctx, err error) error {
	var appErr *apperror.AppError

	// Check if it's an AppError
	if errors.As(err, &appErr) {
		logger.Log.Error().
			Str("code", appErr.Code).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Err(appErr.Err).
			Msg(appErr.Message)

		return Error(c, appErr.StatusCode, appErr.Code, appErr.Message, nil)
	}

	// Unknown error - log and return generic error
	logger.Log.Error().
		Str("path", c.Path()).
		Str("method", c.Method()).
		Err(err).
		Msg("Unhandled error")

	return Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", nil)
}
