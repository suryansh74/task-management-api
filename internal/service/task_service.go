package service

import (
	"context"
	"time"

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

// CreateSession it sets new session
// =========================================================================
func (s *taskService) GetTasks(ctx context.Context) ([]*models.Task, error) {
	// create random id
	tasks := []*models.Task{
		{ID: "1", Title: "Coding in Golang", Content: "very easy", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", Title: "Coding in Golang", Content: "very easy", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "3", Title: "C++ for beginner", Content: "created by bjarne Stroutstrup", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	return tasks, nil
}
