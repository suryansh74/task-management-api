package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/http/response"
	"github.com/suryansh74/task-management-api-project/internal/logger"
)

// ErrorHandler is the global Fiber error handler
// ==================================================
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
// ==================================================
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

// AuthMiddleware checks if incoming req have cookie with valid session user id
// ==================================================
func (s *server) AuthMiddleware(c *fiber.Ctx) error {
	reqCtx := c.UserContext()
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "not logged in",
		})
	}

	userID, err := s.redisClient.HGet(reqCtx, sessionID, "user_id").Result()
	if err != nil || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": "invalid session",
		})
	}

	c.Locals("user_id", userID)
	return c.Next()
}

// RedisRateLimiter
// ==================================================
func (s *server) RedisRateLimiter(
	prefix string,
	limit int,
	window time.Duration,
	keyFunc func(ctx *fiber.Ctx) string,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// incr count
		reqCtx := ctx.UserContext()
		key := fmt.Sprintf("rate:%s:%s", prefix, keyFunc(ctx))
		count, err := s.redisClient.Incr(reqCtx, key).Result()
		if err != nil {
			panic(err)
		}

		// if count == 1 set expiry
		if count == 1 {
			s.redisClient.Expire(reqCtx, key, window)
		}
		// if count > limit block req other wise next
		if count > int64(limit) {
			ttl, _ := s.redisClient.TTL(reqCtx, key).Result()
			return ctx.Status(fiber.StatusTooManyRequests).JSON(&fiber.Map{
				"error":       "too many requests",
				"retry_after": ttl.Seconds(),
			})
		}
		return ctx.Next()
	}
}
