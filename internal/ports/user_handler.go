package ports

import "github.com/gofiber/fiber/v2"

// UserHandler defines the HTTP adapter contract for user/auth operations.
type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}
