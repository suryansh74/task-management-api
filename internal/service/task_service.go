package service

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type taskService struct {
	taskRepo ports.TaskRepository
}

// NewTaskService creates a new user session service instance
// =========================================================================
func NewTaskService(taskRepo ports.TaskRepository) ports.TaskService {
	return &taskService{
		taskRepo: taskRepo,
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
	return s.taskRepo.GetTaskByID(ctx, id)
}

func (s *taskService) UpdateTaskByID(ctx context.Context, id string, task *models.Task) error {
	return s.taskRepo.UpdateTaskByID(ctx, id, task)
}

func (s *taskService) DeleteTaskByID(ctx context.Context, id string) error {
	return s.taskRepo.DeleteTaskByID(ctx, id)
}
