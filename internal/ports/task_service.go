package ports

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/models"
)

// TaskService defines business logic operations for task
type TaskService interface {
	GetTasks(ctx context.Context) ([]*models.Task, error)
}
