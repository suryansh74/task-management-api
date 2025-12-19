package server

import "github.com/gofiber/fiber/v2"

func (s *server) setupRoutes() {
	s.app.Get("/check_health", s.checkHealth)
}

// checkHealth
// ==================================================
func (s *server) checkHealth(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "working fine",
	})
}
