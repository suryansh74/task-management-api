package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type taskRepository struct {
	db *pgx.Conn
}

func NewTaskRepository(db *pgx.Conn) ports.TaskRepository {
	logger.Log.Info().Msg("initializing task repository")
	return &taskRepository{db: db}
}

// GetAllTasks get all tasks
// =========================================================================
func (tr *taskRepository) GetAllTasks(ctx context.Context, userID string) ([]*models.Task, error) {
	logger.Log.Debug().
		Str("user_id", userID).
		Msg("fetching all tasks for user")

	var tasks []*models.Task
	rows, err := tr.db.Query(context.Background(), "select id, title, content, created_at, updated_at from tasks where user_id = $1", userID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to query tasks")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			logger.Log.Error().
				Err(err).
				Str("user_id", userID).
				Msg("failed to scan task row")
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	logger.Log.Info().
		Str("user_id", userID).
		Int("task_count", len(tasks)).
		Msg("successfully fetched all tasks for user")
	return tasks, nil
}

// CreateTask create a task
// =========================================================================
func (tr *taskRepository) CreateTask(ctx context.Context, task *models.Task) (string, error) {
	logger.Log.Debug().
		Str("title", task.Title).
		Str("user_id", task.UserID).
		Msg("creating new task")

	var id string
	err := tr.db.QueryRow(context.Background(), "insert into tasks(title, content, user_id) values($1,$2,$3) returning id", task.Title, task.Content, task.UserID).Scan(&id)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("title", task.Title).
			Str("user_id", task.UserID).
			Msg("failed to create task")
		return "", err
	}

	logger.Log.Info().
		Str("task_id", id).
		Str("title", task.Title).
		Str("user_id", task.UserID).
		Msg("task created successfully")
	return id, nil
}

// Get task by id
// =========================================================================
func (tr *taskRepository) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	logger.Log.Debug().
		Str("task_id", id).
		Msg("fetching task by id")

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
			logger.Log.Warn().
				Str("task_id", id).
				Msg("task not found")
			return nil, apperror.NewNotFoundError("task not found")
		}
		logger.Log.Error().
			Err(err).
			Str("task_id", id).
			Msg("failed to fetch task")
		return nil, err
	}

	logger.Log.Info().
		Str("task_id", task.ID).
		Str("title", task.Title).
		Str("user_id", task.UserID).
		Msg("task fetched successfully")
	return task, nil
}

// Update task by id
// =========================================================================
func (tr *taskRepository) UpdateTaskByID(ctx context.Context, id string, task *models.Task) error {
	logger.Log.Debug().
		Str("task_id", id).
		Str("title", task.Title).
		Msg("updating task")

	cmd, err := tr.db.Exec(ctx,
		`UPDATE tasks
		 SET title = $1, content = $2, updated_at = NOW()
		 WHERE id = $3`,
		task.Title,
		task.Content,
		id,
	)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("task_id", id).
			Msg("failed to update task")
		return err
	}

	if cmd.RowsAffected() == 0 {
		logger.Log.Warn().
			Str("task_id", id).
			Msg("task not found for update")
		return apperror.NewNotFoundError("task not found")
	}

	logger.Log.Info().
		Str("task_id", id).
		Str("title", task.Title).
		Int64("rows_affected", cmd.RowsAffected()).
		Msg("task updated successfully")
	return nil
}

// Delete task by id
// =========================================================================
func (tr *taskRepository) DeleteTaskByID(ctx context.Context, id string) error {
	logger.Log.Debug().
		Str("task_id", id).
		Msg("deleting task")

	cmd, err := tr.db.Exec(ctx,
		`DELETE FROM tasks WHERE id = $1`,
		id,
	)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("task_id", id).
			Msg("failed to delete task")
		return err
	}

	if cmd.RowsAffected() == 0 {
		logger.Log.Warn().
			Str("task_id", id).
			Msg("task not found for deletion")
		return apperror.NewNotFoundError("task not found")
	}

	logger.Log.Info().
		Str("task_id", id).
		Int64("rows_affected", cmd.RowsAffected()).
		Msg("task deleted successfully")
	return nil
}
