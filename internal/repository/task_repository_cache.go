package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type taskCacheRepository struct {
	redisClient *redis.Client
}

func NewTaskCacheRepository(redisClient *redis.Client) ports.TaskCacheRepository {
	logger.Log.Info().Msg("initializing task cache repository")
	return &taskCacheRepository{redisClient: redisClient}
}

// SetTask set task
// =========================================================================
func (s *taskCacheRepository) SetTask(ctx context.Context, task *models.Task, key string, exp time.Duration) error {
	logger.Log.Debug().
		Str("cache_key", key).
		Str("task_id", task.ID).
		Dur("expiration", exp).
		Msg("setting task in cache")

	bytes, err := json.Marshal(task)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("cache_key", key).
			Str("task_id", task.ID).
			Msg("failed to marshal task for caching")
		return err
	}

	err = s.redisClient.Set(ctx, key, bytes, exp).Err()
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("cache_key", key).
			Str("task_id", task.ID).
			Msg("failed to set task in cache")
		return err
	}

	logger.Log.Info().
		Str("cache_key", key).
		Str("task_id", task.ID).
		Dur("expiration", exp).
		Msg("task cached successfully")
	return nil
}

func (s *taskCacheRepository) GetTaskByID(ctx context.Context, key string) (*models.Task, error) {
	logger.Log.Debug().
		Str("cache_key", key).
		Msg("retrieving task from cache")

	val, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.Log.Debug().
			Str("cache_key", key).
			Msg("cache miss: task not found in cache")
		return nil, nil // cache miss
	}
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("cache_key", key).
			Msg("failed to get task from cache")
		return nil, err
	}

	var task models.Task
	if err := json.Unmarshal([]byte(val), &task); err != nil {
		logger.Log.Error().
			Err(err).
			Str("cache_key", key).
			Msg("failed to unmarshal cached task")
		return nil, err
	}

	logger.Log.Info().
		Str("cache_key", key).
		Str("task_id", task.ID).
		Msg("cache hit: task retrieved successfully")
	return &task, nil
}

func (s *taskCacheRepository) DeleteTaskByID(ctx context.Context, key string) error {
	logger.Log.Debug().
		Str("cache_key", key).
		Msg("deleting task from cache")

	err := s.redisClient.Unlink(ctx, key).Err()
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("cache_key", key).
			Msg("failed to delete task from cache")
		return err
	}

	logger.Log.Info().
		Str("cache_key", key).
		Msg("task removed from cache successfully")
	return nil
}
