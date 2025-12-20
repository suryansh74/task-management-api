package ports

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/models"
)

// TaskService defines business logic operations for task
type TaskService interface {
	GetTasks(ctx context.Context) ([]*models.Task, error)
	GetTaskByID(ctx context.Context, id string) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) (string, error)
	UpdateTaskByID(ctx context.Context, id string, task *models.Task) error
	DeleteTaskByID(ctx context.Context, id string) error
}
