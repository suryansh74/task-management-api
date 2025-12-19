package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/handler"
)

// setupRoutes serves all http routes
// ==================================================

func (s *server) setupRoutes(userHandler *handler.UserHandler, taskHandler *handler.TaskHandler) {
	publicLimiter := s.RedisRateLimiter("public", 10, time.Minute, func(c *fiber.Ctx) string {
		return c.IP()
	})
	taskLimiter := s.RedisRateLimiter("task", 100, time.Minute, func(c *fiber.Ctx) string {
		return c.IP()
	})

	// Public route (no auth check)
	s.app.Get("/check_health", publicLimiter, s.checkHealth)

	// Guest-only routes (must NOT be logged in)
	s.app.Post("/register", publicLimiter, s.GuestMiddleware, userHandler.Register)
	s.app.Post("/login", publicLimiter, s.GuestMiddleware, userHandler.Login)

	// Protected routes (must be logged in)
	s.app.Post("/logout", publicLimiter, s.AuthMiddleware, userHandler.Logout)
	s.app.Get("/tasks", taskLimiter, s.AuthMiddleware, taskHandler.GetTasks)
}

// checkHealth
// ==================================================
func (s *server) checkHealth(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "working fine",
	})
}
