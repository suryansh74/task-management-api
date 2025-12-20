package ports

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/models"
)

// TaskService defines business logic operations for task
type TaskService interface {
	GetTasks(ctx context.Context, userID string) ([]*models.Task, error)
	GetTaskByID(ctx context.Context, taskID string, userID string) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) (string, error)
	UpdateTaskByID(ctx context.Context, taskID string, userID string, task *models.Task) error
	DeleteTaskByID(ctx context.Context, taskID string, userID string) error
}
