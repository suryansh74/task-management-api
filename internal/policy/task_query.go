package policy

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type TaskQuery struct {
	cache ports.TaskCacheRepository
	repo  ports.TaskRepository
}

func NewTaskQuery(cache ports.TaskCacheRepository, repo ports.TaskRepository) *TaskQuery {
	return &TaskQuery{cache: cache, repo: repo}
}

func (q *TaskQuery) GetOwnerIDByTaskID(ctx context.Context, taskID string, key string) (string, error) {
	// 1️⃣ Try cache
	task, _ := q.cache.GetTaskByID(ctx, key)
	if task != nil {
		return task.UserID, nil
	}

	// 2️⃣ Fallback to DB
	task, err := q.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return "", err
	}

	return task.UserID, nil
}
