package service

import (
	"context"
	"fmt"
	"time"

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
func (s *taskService) GetTasks(ctx context.Context) ([]*models.Task, error) {
	tasks, err := s.taskRepo.GetAllTasks(ctx)
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

func (s *taskService) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
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

func (s *taskService) UpdateTaskByID(ctx context.Context, id string, task *models.Task) error {
	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, id)
	err := s.taskRepo.UpdateTaskByID(ctx, id, task)
	if err != nil {
		return err
	}
	s.taskCacheRepo.DeleteTaskByID(ctx, key)
	return nil
}

func (s *taskService) DeleteTaskByID(ctx context.Context, id string) error {
	key := fmt.Sprintf("%s:cache:task:%s", s.redisAppName, id)
	err := s.taskRepo.DeleteTaskByID(ctx, id)
	if err != nil {
		return err
	}
	s.taskCacheRepo.DeleteTaskByID(ctx, key)
	return nil
}
