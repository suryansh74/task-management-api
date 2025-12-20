package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/http/response"
	"github.com/suryansh74/task-management-api-project/internal/logger"
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
	logger.Log.Info().
		Str("redis_app_name", redisAppName).
		Dur("cache_expiration", cacheExpiration).
		Msg("initializing task handler")
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
	logger.Log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("ip", c.IP()).
		Msg("received request to get all tasks")

	// policy
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		logger.Log.Warn().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Msg("unauthorized request: invalid auth context")
		return apperror.NewUnauthorizedError("invalid auth context")
	}

	logger.Log.Debug().
		Str("user_id", userID).
		Msg("authenticated user fetching tasks")

	// Call service
	tasks, err := h.taskService.GetTasks(c.Context(), userID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("user_id", userID).
			Str("path", c.Path()).
			Msg("failed to fetch tasks")
		return err // Global error handler will catch this
	}

	logger.Log.Info().
		Str("user_id", userID).
		Int("task_count", len(tasks)).
		Int("status", fiber.StatusOK).
		Msg("successfully returned all tasks")

	return response.Success(c, fiber.StatusOK, "All Returned Tasks", tasks)
}

// CreateTask create task
// =========================================================================
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	logger.Log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("ip", c.IP()).
		Msg("received request to create task")

	var req models.Task

	// Parse body
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Warn().
			Err(err).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Msg("failed to parse request body")
		return apperror.NewBadRequestError("Invalid request body")
	}

	logger.Log.Debug().
		Str("title", req.Title).
		Msg("parsed task creation request")

	// Validate
	if fieldErrors := validator.ValidateStruct(req); len(fieldErrors) > 0 {
		logger.Log.Warn().
			Interface("validation_errors", fieldErrors).
			Str("title", req.Title).
			Msg("validation failed for task creation")
		return response.ValidationError(c, fieldErrors)
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		logger.Log.Warn().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Msg("unauthorized request: invalid auth context")
		return apperror.NewUnauthorizedError("invalid auth context")
	}

	req.UserID = userID

	logger.Log.Debug().
		Str("user_id", userID).
		Str("title", req.Title).
		Msg("creating task for user")

	// Call service
	id, err := h.taskService.CreateTask(c.Context(), &req)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("user_id", userID).
			Str("title", req.Title).
			Msg("failed to create task")
		return err // Global error handler will catch this
	}

	logger.Log.Info().
		Str("task_id", id).
		Str("user_id", userID).
		Str("title", req.Title).
		Int("status", fiber.StatusOK).
		Msg("task created successfully")

	return response.Success(c, fiber.StatusOK, "Task Created", id)
}

// GetTaskByID get task
// =========================================================================
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")

	logger.Log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("task_id", id).
		Str("ip", c.IP()).
		Msg("received request to get task by id")

	if id == "" {
		logger.Log.Warn().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Msg("missing task id in request")
		return apperror.NewBadRequestError("task id is required")
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		logger.Log.Warn().
			Str("task_id", id).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Msg("unauthorized request: invalid auth context")
		return apperror.NewUnauthorizedError("invalid auth context")
	}

	logger.Log.Debug().
		Str("task_id", id).
		Str("user_id", userID).
		Msg("fetching task for user")

	task, err := h.taskService.GetTaskByID(c.Context(), id, userID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("task_id", id).
			Str("user_id", userID).
			Msg("failed to get task")
		return err
	}

	logger.Log.Info().
		Str("task_id", id).
		Str("user_id", userID).
		Int("status", fiber.StatusOK).
		Msg("task retrieved successfully")

	return response.Success(c, fiber.StatusOK, "Task Found", task)
}

// UpdateTaskByID update task
// =========================================================================
func (h *TaskHandler) UpdateTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")

	logger.Log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("task_id", id).
		Str("ip", c.IP()).
		Msg("received request to update task")

	if id == "" {
		logger.Log.Warn().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Msg("missing task id in request")
		return apperror.NewBadRequestError("task id is required")
	}

	var req UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Warn().
			Err(err).
			Str("task_id", id).
			Str("method", c.Method()).
			Msg("failed to parse request body")
		return apperror.NewBadRequestError("invalid request body")
	}

	logger.Log.Debug().
		Str("task_id", id).
		Str("title", req.Title).
		Msg("parsed task update request")

	if fieldErrors := validator.ValidateStruct(req); len(fieldErrors) > 0 {
		logger.Log.Warn().
			Interface("validation_errors", fieldErrors).
			Str("task_id", id).
			Str("title", req.Title).
			Msg("validation failed for task update")
		return response.ValidationError(c, fieldErrors)
	}

	task := &models.Task{
		Title:   req.Title,
		Content: req.Content,
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		logger.Log.Warn().
			Str("task_id", id).
			Str("method", c.Method()).
			Msg("unauthorized request: invalid auth context")
		return apperror.NewUnauthorizedError("invalid auth context")
	}

	logger.Log.Debug().
		Str("task_id", id).
		Str("user_id", userID).
		Str("title", req.Title).
		Msg("updating task for user")

	if err := h.taskService.UpdateTaskByID(c.Context(), id, userID, task); err != nil {
		logger.Log.Error().
			Err(err).
			Str("task_id", id).
			Str("user_id", userID).
			Msg("failed to update task")
		return err
	}

	logger.Log.Info().
		Str("task_id", id).
		Str("user_id", userID).
		Str("title", req.Title).
		Int("status", fiber.StatusOK).
		Msg("task updated successfully")

	return response.Success(c, fiber.StatusOK, "Task Updated", nil)
}

// DeleteTaskByID delete task
// =========================================================================
func (h *TaskHandler) DeleteTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")

	logger.Log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("task_id", id).
		Str("ip", c.IP()).
		Msg("received request to delete task")

	if id == "" {
		logger.Log.Warn().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Msg("missing task id in request")
		return apperror.NewBadRequestError("task id is required")
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		logger.Log.Warn().
			Str("task_id", id).
			Str("method", c.Method()).
			Msg("unauthorized request: invalid auth context")
		return apperror.NewUnauthorizedError("invalid auth context")
	}

	logger.Log.Debug().
		Str("task_id", id).
		Str("user_id", userID).
		Msg("deleting task for user")

	if err := h.taskService.DeleteTaskByID(c.Context(), id, userID); err != nil {
		logger.Log.Error().
			Err(err).
			Str("task_id", id).
			Str("user_id", userID).
			Msg("failed to delete task")
		return err
	}

	logger.Log.Info().
		Str("task_id", id).
		Str("user_id", userID).
		Int("status", fiber.StatusOK).
		Msg("task deleted successfully")

	return response.Success(c, fiber.StatusOK, "Task Deleted", nil)
}
