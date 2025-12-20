package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/http/response"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
	"github.com/suryansh74/task-management-api-project/internal/validator"
)

type TaskHandler struct {
	taskService     ports.TaskService
	redisAppName    string
	cacheExpiration time.Duration
}

// NewTaskHandler Constructor for TaskHandler
// =========================================================================
func NewTaskHandler(taskService ports.TaskService, redisAppName string, cacheExpiration time.Duration) *TaskHandler {
	return &TaskHandler{
		taskService:     taskService,
		redisAppName:    redisAppName,
		cacheExpiration: cacheExpiration,
	}
}

// CreateTaskRequest dto for incoming req
// =========================================================================
type CreateTaskRequest struct {
	Title   string `json:"title" validate:"required,min=2,max=100"`
	Content string `json:"content"`
}

type UpdateTaskRequest struct {
	Title   string `json:"title" validate:"min=2,max=100"`
	Content string `json:"content" validate:"max=500"`
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

// CreateTask create task
// =========================================================================
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var req models.Task

	// Parse body
	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequestError("Invalid request body")
	}

	// Validate
	if fieldErrors := validator.ValidateStruct(req); len(fieldErrors) > 0 {
		return response.ValidationError(c, fieldErrors)
	}

	// Call service
	id, err := h.taskService.CreateTask(c.Context(), &req)
	if err != nil {
		return err // Global error handler will catch this
	}

	return response.Success(c, fiber.StatusOK, "Task Created", id)
}

// GetTaskByID get task
// =========================================================================
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperror.NewBadRequestError("task id is required")
	}

	task, err := h.taskService.GetTaskByID(c.Context(), id)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Task Found", task)
}

// UpdateTaskByID update task
// =========================================================================
func (h *TaskHandler) UpdateTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperror.NewBadRequestError("task id is required")
	}

	var req UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequestError("invalid request body")
	}

	if fieldErrors := validator.ValidateStruct(req); len(fieldErrors) > 0 {
		return response.ValidationError(c, fieldErrors)
	}

	task := &models.Task{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.taskService.UpdateTaskByID(c.Context(), id, task); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Task Updated", nil)
}

// DeleteTaskByID delete task
// =========================================================================
func (h *TaskHandler) DeleteTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperror.NewBadRequestError("task id is required")
	}

	if err := h.taskService.DeleteTaskByID(c.Context(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Task Deleted", nil)
}
