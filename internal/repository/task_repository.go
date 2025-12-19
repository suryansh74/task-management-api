package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type taskRepository struct {
	db *pgx.Conn
}

func NewTaskRepository(db *pgx.Conn) ports.TaskRepository {
	return &taskRepository{db: db}
}

func (tr *taskRepository) GetAllTasks(ctx context.Context) ([]*models.Task, error) {
	return nil, nil
}
