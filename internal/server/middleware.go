package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/http/response"
	"github.com/suryansh74/task-management-api-project/internal/logger"
)

// ErrorHandler is the global Fiber error handler
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Handle Fiber's own errors
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return response.Error(c, fiberErr.Code, "REQUEST_ERROR", fiberErr.Message, nil)
		}

		// Handle application errors
		return response.HandleError(c, err)
	}
}

// RequestLogger logs incoming requests
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Msg("Incoming request")

		return c.Next()
	}
}
