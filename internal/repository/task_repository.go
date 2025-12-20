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

// GetAllTasks get all tasks
// =========================================================================
func (tr *taskRepository) GetAllTasks(ctx context.Context) ([]*models.Task, error) {
	var tasks []*models.Task
	rows, _ := tr.db.Query(context.Background(), "select * from tasks")

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
	err := tr.db.QueryRow(context.Background(), "insert into tasks(title, content) values($1,$2) returning id", task.Title, task.Content).Scan(&id)
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
		`SELECT id, title, content, created_at, updated_at
		 FROM tasks WHERE id = $1`,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Content,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
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
		return pgx.ErrNoRows
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
		return pgx.ErrNoRows
	}

	return nil
}
