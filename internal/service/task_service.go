package service

import (
	"context"
	"fmt"
	"time"

	"github.com/suryansh74/task-management-api-project/internal/apperror"
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
	tasks, err := s.taskRepo.GetAllTasks(ctx, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// CreateTask get all tasks
// =========================================================================
func (s *taskService) CreateTask(ctx context.Context, task *models.Task) (string, error) {
	id, err := s.taskRepo.CreateTask(ctx, task)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetTaskByID get tasks
// =========================================================================
func (s *taskService) GetTaskByID(ctx context.Context, taskID string, userID string) (*models.Task, error) {
	// check policy
	_, err := s.mustBeOwner(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}
	// first get from cache
	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, taskID)
	task, _ := s.taskCacheRepo.GetTaskByID(ctx, key)
	if task != nil {
		return task, nil
	}

	// if not exist then from db
	task, err = s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// set in cache
	s.taskCacheRepo.SetTask(ctx, task, key, s.cacheExpiration)
	return task, nil
}

// UpdateTaskByID update tasks
// =========================================================================
func (s *taskService) UpdateTaskByID(ctx context.Context, taskID string, userID string, task *models.Task) error {
	// check policy
	_, err := s.mustBeOwner(ctx, userID, taskID)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, taskID)
	err = s.taskRepo.UpdateTaskByID(ctx, taskID, task)
	if err != nil {
		return err
	}
	// invalidate cache
	s.taskCacheRepo.DeleteTaskByID(ctx, key)
	return nil
}

// DeleteTaskByID delete tasks
// =========================================================================
func (s *taskService) DeleteTaskByID(ctx context.Context, taskID string, userID string) error {
	// check policy
	_, err := s.mustBeOwner(ctx, userID, taskID)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, taskID)
	err = s.taskRepo.DeleteTaskByID(ctx, taskID)
	if err != nil {
		return err
	}
	s.taskCacheRepo.DeleteTaskByID(ctx, key)
	return nil
}

// mustBeOwner helper function to check ownership
// =========================================================================
func (s *taskService) mustBeOwner(
	ctx context.Context,
	userID, taskID string,
) (*models.Task, error) {
	if userID == "" {
		return nil, apperror.NewUnauthorizedError("not authenticated")
	}

	task, err := s.getTaskByIDHelper(ctx, taskID)
	if err != nil {
		return nil, err // not found bubbles up
	}

	if task.UserID != userID {
		return nil, apperror.NewForbiddenError("not allowed")
	}

	return task, nil
}

// getTaskByIDHelper helper function to get task by id without checking ownership
// =========================================================================
func (s *taskService) getTaskByIDHelper(ctx context.Context, id string) (*models.Task, error) {
	// first get from cache
	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, id)
	task, _ := s.taskCacheRepo.GetTaskByID(ctx, key)
	if task != nil {
		return task, nil
	}

	// if not exist then from db
	task, err := s.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// set in cache
	s.taskCacheRepo.SetTask(ctx, task, key, s.cacheExpiration)
	return task, nil
}
