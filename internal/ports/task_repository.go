package ports

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/models"
)

type TaskRepository interface {
	GetAllTasks(ctx context.Context) ([]*models.Task, error)
}
