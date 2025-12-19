package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/handler"
)

// setupRoutes serves all http routes
// ==================================================
func (s *server) setupRoutes(userHandler *handler.UserHandler, taskHandler *handler.TaskHandler) {
	// init redis rate limit
	publicLimiter := s.RedisRateLimiter("public", 10, time.Minute, func(c *fiber.Ctx) string {
		return c.IP()
	})
	taskLimiter := s.RedisRateLimiter("task", 100, time.Minute, func(c *fiber.Ctx) string {
		return c.IP()
	})
	// public routes
	s.app.Get("/check_health", publicLimiter, s.checkHealth)
	s.app.Post("/register", publicLimiter, userHandler.Register)
	s.app.Post("/login", publicLimiter, userHandler.Login)

	// protected routes
	protected := s.app.Group("/", s.AuthMiddleware)
	protected.Get("/tasks", taskLimiter, taskHandler.GetTasks)
}

// checkHealth
// ==================================================
func (s *server) checkHealth(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "working fine",
	})
}
