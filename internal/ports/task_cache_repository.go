package ports

import (
	"context"
	"time"

	"github.com/suryansh74/task-management-api-project/internal/models"
)

type TaskCacheRepository interface {
	SetTask(
		ctx context.Context,
		task *models.Task,
		key string,
		expiration time.Duration,
	) error

	GetTaskByID(
		ctx context.Context,
		key string,
	) (*models.Task, error)

	DeleteTaskByID(
		ctx context.Context,
		key string,
	) error
}
