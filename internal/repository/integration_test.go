//go:build integration

package repository_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/repository"
)

func init() {
	logger.Init()
}

func findInitSQL(t *testing.T) string {
	t.Helper()
	candidates := []string{
		"init.sql",
		"../../init.sql",
		"../init.sql",
	}
	// Also try from module root via env or walk
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, "init.sql"),
			filepath.Join(wd, "..", "..", "init.sql"),
		)
	}
	for _, c := range candidates {
		if abs, err := filepath.Abs(c); err == nil {
			if _, err := os.Stat(abs); err == nil {
				return abs
			}
		}
	}
	t.Fatal("init.sql not found")
	return ""
}

func setupPostgres(t *testing.T) (*pgx.Conn, func()) {
	t.Helper()
	ctx := context.Background()
	initSQL := findInitSQL(t)

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("task_management_api"),
		tcpostgres.WithUsername("root"),
		tcpostgres.WithPassword("secret"),
		tcpostgres.WithInitScripts(initSQL),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	conn, err := pgx.Connect(ctx, connStr)
	require.NoError(t, err)

	cleanup := func() {
		conn.Close(ctx)
		_ = pgContainer.Terminate(ctx)
	}
	return conn, cleanup
}

func setupRedis(t *testing.T) (*goredis.Client, func()) {
	t.Helper()
	ctx := context.Background()

	redisContainer, err := tcredis.Run(ctx, "redis:7-alpine")
	require.NoError(t, err)

	endpoint, err := redisContainer.Endpoint(ctx, "")
	require.NoError(t, err)

	client := goredis.NewClient(&goredis.Options{Addr: endpoint})
	require.NoError(t, client.Ping(ctx).Err())

	cleanup := func() {
		_ = client.Close()
		_ = redisContainer.Terminate(ctx)
	}
	return client, cleanup
}

func TestUserRepository_Integration(t *testing.T) {
	conn, cleanup := setupPostgres(t)
	defer cleanup()

	repo := repository.NewUserRepository(conn)
	ctx := context.Background()

	t.Run("CreateUser and FindByEmail", func(t *testing.T) {
		user := &models.User{
			Name:     "Integration User",
			Email:    "integration@example.com",
			Password: "hashed-password-123",
		}

		id, err := repo.CreateUser(ctx, user)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		found, err := repo.FindByEmail(ctx, "integration@example.com")
		require.NoError(t, err)
		require.Equal(t, id, found.ID)
		require.Equal(t, "Integration User", found.Name)
		require.Equal(t, "integration@example.com", found.Email)
	})

	t.Run("CreateUser duplicate email returns conflict", func(t *testing.T) {
		user := &models.User{
			Name:     "Dup User",
			Email:    "dup@example.com",
			Password: "pass",
		}
		_, err := repo.CreateUser(ctx, user)
		require.NoError(t, err)

		_, err = repo.CreateUser(ctx, user)
		require.Error(t, err)
	})

	t.Run("FindByEmail not found", func(t *testing.T) {
		_, err := repo.FindByEmail(ctx, "does-not-exist@example.com")
		require.Error(t, err)
	})
}

func TestTaskRepository_Integration(t *testing.T) {
	conn, cleanup := setupPostgres(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(conn)
	taskRepo := repository.NewTaskRepository(conn)
	ctx := context.Background()

	userID, err := userRepo.CreateUser(ctx, &models.User{
		Name:     "Task Owner",
		Email:    "taskowner@example.com",
		Password: "pass",
	})
	require.NoError(t, err)

	t.Run("CRUD lifecycle", func(t *testing.T) {
		task := &models.Task{
			UserID:  userID,
			Title:   "Integration Task",
			Content: "Created by testcontainers",
		}
		taskID, err := taskRepo.CreateTask(ctx, task)
		require.NoError(t, err)
		require.NotEmpty(t, taskID)

		got, err := taskRepo.GetTaskByID(ctx, taskID)
		require.NoError(t, err)
		require.Equal(t, "Integration Task", got.Title)
		require.Equal(t, "Created by testcontainers", got.Content)

		all, err := taskRepo.GetAllTasks(ctx, userID)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(all), 1)

		err = taskRepo.UpdateTaskByID(ctx, taskID, &models.Task{
			Title:   "Updated Title",
			Content: "Updated content",
		})
		require.NoError(t, err)

		updated, err := taskRepo.GetTaskByID(ctx, taskID)
		require.NoError(t, err)
		require.Equal(t, "Updated Title", updated.Title)

		err = taskRepo.DeleteTaskByID(ctx, taskID)
		require.NoError(t, err)

		_, err = taskRepo.GetTaskByID(ctx, taskID)
		require.Error(t, err)
	})

	t.Run("GetAllTasks empty for unknown user", func(t *testing.T) {
		tasks, err := taskRepo.GetAllTasks(ctx, "00000000-0000-0000-0000-000000000000")
		require.NoError(t, err)
		require.Empty(t, tasks)
	})
}

func TestTaskCacheRepository_Integration(t *testing.T) {
	rdb, cleanup := setupRedis(t)
	defer cleanup()

	cacheRepo := repository.NewTaskCacheRepository(rdb)
	ctx := context.Background()

	task := &models.Task{
		ID:        "cache-task-1",
		UserID:    "user-1",
		Title:     "Cached Task",
		Content:   "from redis",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	key := "task:cache-task-1"

	t.Run("Set Get Delete", func(t *testing.T) {
		err := cacheRepo.SetTask(ctx, task, key, 2*time.Minute)
		require.NoError(t, err)

		got, err := cacheRepo.GetTaskByID(ctx, key)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, task.ID, got.ID)
		require.Equal(t, task.Title, got.Title)

		err = cacheRepo.DeleteTaskByID(ctx, key)
		require.NoError(t, err)

		// cache miss returns (nil, nil) in this implementation
		got, err = cacheRepo.GetTaskByID(ctx, key)
		require.NoError(t, err)
		require.Nil(t, got)
	})
}
