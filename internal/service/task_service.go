package service

import (
	"context"
	"fmt"
	"time"

	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type taskService struct {
	taskRepo        ports.TaskRepository
	taskCacheRepo   ports.TaskCacheRepository
	redisAppName    string
	cacheExpiration time.Duration
}

// NewTaskService creates a new user session service instance
// =========================================================================
func NewTaskService(taskRepo ports.TaskRepository, taskCacheRepo ports.TaskCacheRepository, redisAppName string, cacheExpiration time.Duration) ports.TaskService {
	logger.Log.Info().
		Str("redis_app_name", redisAppName).
		Dur("cache_expiration", cacheExpiration).
		Msg("initializing task service")
	return &taskService{
		taskRepo:        taskRepo,
		taskCacheRepo:   taskCacheRepo,
		redisAppName:    redisAppName,
		cacheExpiration: cacheExpiration,
	}
}

// GetTasks get all tasks
// =========================================================================
func (s *taskService) GetTasks(ctx context.Context, userID string) ([]*models.Task, error) {
	logger.Log.Debug().
		Str("user_id", userID).
		Msg("fetching all tasks for user")

	tasks, err := s.taskRepo.GetAllTasks(ctx, userID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to fetch tasks")
		return nil, err
	}

	logger.Log.Info().
		Str("user_id", userID).
		Int("task_count", len(tasks)).
		Msg("tasks fetched successfully")
	return tasks, nil
}

// CreateTask get all tasks
// =========================================================================
func (s *taskService) CreateTask(ctx context.Context, task *models.Task) (string, error) {
	logger.Log.Debug().
		Str("user_id", task.UserID).
		Str("title", task.Title).
		Msg("creating new task")

	id, err := s.taskRepo.CreateTask(ctx, task)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("user_id", task.UserID).
			Str("title", task.Title).
			Msg("failed to create task")
		return "", err
	}

	logger.Log.Info().
		Str("task_id", id).
		Str("user_id", task.UserID).
		Str("title", task.Title).
		Msg("task created successfully")
	return id, nil
}

// GetTaskByID get tasks
// =========================================================================
func (s *taskService) GetTaskByID(ctx context.Context, taskID string, userID string) (*models.Task, error) {
	logger.Log.Debug().
		Str("task_id", taskID).
		Str("user_id", userID).
		Msg("fetching task by id")

	// check policy
	_, err := s.mustBeOwner(ctx, userID, taskID)
	if err != nil {
		logger.Log.Warn().
			Err(err).
			Str("task_id", taskID).
			Str("user_id", userID).
			Msg("ownership validation failed")
		return nil, err
	}

	// first get from cache
	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, taskID)
	logger.Log.Debug().
		Str("cache_key", key).
		Str("task_id", taskID).
		Msg("checking cache for task")

	task, _ := s.taskCacheRepo.GetTaskByID(ctx, key)
	if task != nil {
		logger.Log.Info().
			Str("task_id", taskID).
			Str("user_id", userID).
			Str("cache_key", key).
			Msg("task retrieved from cache")
		return task, nil
	}

	// if not exist then from db
	logger.Log.Debug().
		Str("task_id", taskID).
		Msg("cache miss, fetching from database")

	task, err = s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("task_id", taskID).
			Str("user_id", userID).
			Msg("failed to fetch task from database")
		return nil, err
	}

	// set in cache
	logger.Log.Debug().
		Str("cache_key", key).
		Str("task_id", taskID).
		Msg("caching task for future requests")
	s.taskCacheRepo.SetTask(ctx, task, key, s.cacheExpiration)

	logger.Log.Info().
		Str("task_id", taskID).
		Str("user_id", userID).
		Msg("task retrieved successfully from database")
	return task, nil
}

// UpdateTaskByID update tasks
// =========================================================================
func (s *taskService) UpdateTaskByID(ctx context.Context, taskID string, userID string, task *models.Task) error {
	logger.Log.Debug().
		Str("task_id", taskID).
		Str("user_id", userID).
		Str("title", task.Title).
		Msg("updating task")

	// check policy
	_, err := s.mustBeOwner(ctx, userID, taskID)
	if err != nil {
		logger.Log.Warn().
			Err(err).
			Str("task_id", taskID).
			Str("user_id", userID).
			Msg("ownership validation failed for update")
		return err
	}

	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, taskID)
	err = s.taskRepo.UpdateTaskByID(ctx, taskID, task)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("task_id", taskID).
			Str("user_id", userID).
			Msg("failed to update task")
		return err
	}

	// invalidate cache
	logger.Log.Debug().
		Str("cache_key", key).
		Str("task_id", taskID).
		Msg("invalidating task cache")
	s.taskCacheRepo.DeleteTaskByID(ctx, key)

	logger.Log.Info().
		Str("task_id", taskID).
		Str("user_id", userID).
		Str("title", task.Title).
		Msg("task updated successfully")
	return nil
}

// DeleteTaskByID delete tasks
// =========================================================================
func (s *taskService) DeleteTaskByID(ctx context.Context, taskID string, userID string) error {
	logger.Log.Debug().
		Str("task_id", taskID).
		Str("user_id", userID).
		Msg("deleting task")

	// check policy
	_, err := s.mustBeOwner(ctx, userID, taskID)
	if err != nil {
		logger.Log.Warn().
			Err(err).
			Str("task_id", taskID).
			Str("user_id", userID).
			Msg("ownership validation failed for deletion")
		return err
	}

	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, taskID)
	err = s.taskRepo.DeleteTaskByID(ctx, taskID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("task_id", taskID).
			Str("user_id", userID).
			Msg("failed to delete task")
		return err
	}

	logger.Log.Debug().
		Str("cache_key", key).
		Str("task_id", taskID).
		Msg("removing task from cache")
	s.taskCacheRepo.DeleteTaskByID(ctx, key)

	logger.Log.Info().
		Str("task_id", taskID).
		Str("user_id", userID).
		Msg("task deleted successfully")
	return nil
}

// mustBeOwner helper function to check ownership
// =========================================================================
func (s *taskService) mustBeOwner(
	ctx context.Context,
	userID, taskID string,
) (*models.Task, error) {
	logger.Log.Debug().
		Str("user_id", userID).
		Str("task_id", taskID).
		Msg("validating task ownership")

	if userID == "" {
		logger.Log.Warn().
			Str("task_id", taskID).
			Msg("unauthenticated access attempt")
		return nil, apperror.NewUnauthorizedError("not authenticated")
	}

	task, err := s.getTaskByIDHelper(ctx, taskID)
	if err != nil {
		logger.Log.Warn().
			Err(err).
			Str("user_id", userID).
			Str("task_id", taskID).
			Msg("task not found during ownership check")
		return nil, err // not found bubbles up
	}

	if task.UserID != userID {
		logger.Log.Warn().
			Str("user_id", userID).
			Str("task_id", taskID).
			Str("task_owner_id", task.UserID).
			Msg("unauthorized access attempt: user is not task owner")
		return nil, apperror.NewForbiddenError("not allowed")
	}

	logger.Log.Debug().
		Str("user_id", userID).
		Str("task_id", taskID).
		Msg("ownership validation successful")
	return task, nil
}

// getTaskByIDHelper helper function to get task by id without checking ownership
// =========================================================================
func (s *taskService) getTaskByIDHelper(ctx context.Context, id string) (*models.Task, error) {
	logger.Log.Debug().
		Str("task_id", id).
		Msg("fetching task by id (helper)")

	// first get from cache
	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, id)
	task, _ := s.taskCacheRepo.GetTaskByID(ctx, key)
	if task != nil {
		logger.Log.Debug().
			Str("task_id", id).
			Str("cache_key", key).
			Msg("task found in cache (helper)")
		return task, nil
	}

	// if not exist then from db
	logger.Log.Debug().
		Str("task_id", id).
		Msg("cache miss, fetching from database (helper)")

	task, err := s.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("task_id", id).
			Msg("failed to fetch task from database (helper)")
		return nil, err
	}

	// set in cache
	logger.Log.Debug().
		Str("cache_key", key).
		Str("task_id", id).
		Msg("caching task (helper)")
	s.taskCacheRepo.SetTask(ctx, task, key, s.cacheExpiration)

	logger.Log.Debug().
		Str("task_id", id).
		Msg("task retrieved successfully (helper)")
	return task, nil
}
