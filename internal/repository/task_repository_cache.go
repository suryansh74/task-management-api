package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type taskCacheRepository struct {
	redisClient *redis.Client
}

func NewTaskCacheRepository(redisClient *redis.Client) ports.TaskCacheRepository {
	return &taskCacheRepository{redisClient: redisClient}
}

// SetTask set task
// =========================================================================
func (s *taskCacheRepository) SetTask(ctx context.Context, task *models.Task, key string, exp time.Duration) error {
	bytes, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, bytes, exp).Err()
}

func (s *taskCacheRepository) GetTaskByID(ctx context.Context, key string) (*models.Task, error) {
	val, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	}
	if err != nil {
		return nil, err
	}

	var task models.Task
	if err := json.Unmarshal([]byte(val), &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *taskCacheRepository) DeleteTaskByID(ctx context.Context, key string) error {
	return s.redisClient.Unlink(ctx, key).Err()
}
