package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/handler"
)

func (s *server) setupRoutes(userHandler *handler.UserHandler) {
	s.app.Get("/check_health", s.checkHealth)
	s.app.Post("/users", userHandler.CreateUser)
	s.app.Get("/users", userHandler.GetUser)
}

// checkHealth
// ==================================================
func (s *server) checkHealth(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "working fine",
	})
}
