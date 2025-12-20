package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type taskRepository struct {
	db *pgx.Conn
}

func NewTaskRepository(db *pgx.Conn) ports.TaskRepository {
	return &taskRepository{db: db}
}

// GetAllTasks get all tasks
// =========================================================================
func (tr *taskRepository) GetAllTasks(ctx context.Context, userID string) ([]*models.Task, error) {
	var tasks []*models.Task
	rows, _ := tr.db.Query(context.Background(), "select id, title, content, created_at, updated_at from tasks where user_id = $1", userID)

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

// CreateTask create a task
// =========================================================================
func (tr *taskRepository) CreateTask(ctx context.Context, task *models.Task) (string, error) {
	var id string
	err := tr.db.QueryRow(context.Background(), "insert into tasks(title, content, user_id) values($1,$2,$3) returning id", task.Title, task.Content, task.UserID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Get task by id
// =========================================================================
func (tr *taskRepository) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	task := new(models.Task)

	err := tr.db.QueryRow(ctx,
		`SELECT id, title, content, created_at, updated_at, user_id
		 FROM tasks WHERE id = $1`,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Content,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.UserID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFoundError("task not found")
		}
		return nil, err
	}

	return task, nil
}

// Update task by id
// =========================================================================
func (tr *taskRepository) UpdateTaskByID(ctx context.Context, id string, task *models.Task) error {
	cmd, err := tr.db.Exec(ctx,
		`UPDATE tasks
		 SET title = $1, content = $2, updated_at = NOW()
		 WHERE id = $3`,
		task.Title,
		task.Content,
		id,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return apperror.NewNotFoundError("task not found")
	}

	return nil
}

// Delete task by id
// =========================================================================
func (tr *taskRepository) DeleteTaskByID(ctx context.Context, id string) error {
	cmd, err := tr.db.Exec(ctx,
		`DELETE FROM tasks WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return apperror.NewNotFoundError("task not found")
	}

	return nil
}
