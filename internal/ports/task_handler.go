package ports

import "github.com/gofiber/fiber/v2"

// TaskHandler defines the HTTP adapter contract for task operations.
// Both the concrete REST handler and any future adapters can satisfy this.
type TaskHandler interface {
	GetTasks(c *fiber.Ctx) error
	CreateTask(c *fiber.Ctx) error
	GetTaskByID(c *fiber.Ctx) error
	UpdateTaskByID(c *fiber.Ctx) error
	DeleteTaskByID(c *fiber.Ctx) error
}
