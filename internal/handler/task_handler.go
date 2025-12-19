package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/http/response"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type TaskHandler struct {
	taskService ports.TaskService
}

// NewTaskHandler Constructor for TaskHandler
// =========================================================================
func NewTaskHandler(taskService ports.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTaskRequest dto for incoming req
// =========================================================================
type CreateTaskRequest struct {
	Title   string `json:"title" validate:"required,min=2,max=100"`
	Content string `json:"content"`
}

// GetTaskResponse dto for incoming req
// =========================================================================
type GetTaskResponse struct {
	ID        string    `json:"task_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetTasks return all tasks
// =========================================================================
func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	// Call service
	tasks, err := h.taskService.GetTasks(c.Context())
	if err != nil {
		return err // Global error handler will catch this
	}

	return response.Success(c, fiber.StatusOK, "All Returned Tasks", tasks)
}
