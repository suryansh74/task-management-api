package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

func init() {
	logger.Init()
}

type mockTaskRepository struct {
	getAllFn     func(ctx context.Context, userID string) ([]*models.Task, error)
	getByIDFn    func(ctx context.Context, id string) (*models.Task, error)
	createFn     func(ctx context.Context, task *models.Task) (string, error)
	updateFn     func(ctx context.Context, id string, task *models.Task) error
	deleteFn     func(ctx context.Context, id string) error
}

func (m *mockTaskRepository) GetAllTasks(ctx context.Context, userID string) ([]*models.Task, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, userID)
	}
	return nil, errors.New("not implemented")
}
func (m *mockTaskRepository) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (m *mockTaskRepository) CreateTask(ctx context.Context, task *models.Task) (string, error) {
	if m.createFn != nil {
		return m.createFn(ctx, task)
	}
	return "", errors.New("not implemented")
}
func (m *mockTaskRepository) UpdateTaskByID(ctx context.Context, id string, task *models.Task) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, task)
	}
	return errors.New("not implemented")
}
func (m *mockTaskRepository) DeleteTaskByID(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return errors.New("not implemented")
}

type mockTaskCacheRepository struct {
	getFn    func(ctx context.Context, key string) (*models.Task, error)
	setFn    func(ctx context.Context, task *models.Task, key string, exp time.Duration) error
	deleteFn func(ctx context.Context, key string) error
}

func (m *mockTaskCacheRepository) GetTaskByID(ctx context.Context, key string) (*models.Task, error) {
	if m.getFn != nil {
		return m.getFn(ctx, key)
	}
	return nil, nil
}
func (m *mockTaskCacheRepository) SetTask(ctx context.Context, task *models.Task, key string, exp time.Duration) error {
	if m.setFn != nil {
		return m.setFn(ctx, task, key, exp)
	}
	return nil
}
func (m *mockTaskCacheRepository) DeleteTaskByID(ctx context.Context, key string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, key)
	}
	return nil
}

func TestTaskService_CreateTask(t *testing.T) {
	repo := &mockTaskRepository{
		createFn: func(ctx context.Context, task *models.Task) (string, error) {
			return "task-1", nil
		},
	}
	cache := &mockTaskCacheRepository{}
	svc := NewTaskService(repo, cache, "app", 10*time.Minute)

	id, err := svc.CreateTask(context.Background(), &models.Task{
		UserID:  "user-1",
		Title:   "Test Task",
		Content: "Content",
	})
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if id != "task-1" {
		t.Errorf("expected task-1, got %s", id)
	}
}

func TestTaskService_GetTasks(t *testing.T) {
	expected := []*models.Task{{ID: "t1", Title: "A"}, {ID: "t2", Title: "B"}}
	repo := &mockTaskRepository{
		getAllFn: func(ctx context.Context, userID string) ([]*models.Task, error) {
			if userID != "user-1" {
				t.Errorf("unexpected userID: %s", userID)
			}
			return expected, nil
		},
	}
	cache := &mockTaskCacheRepository{}
	svc := NewTaskService(repo, cache, "app", 10*time.Minute)

	tasks, err := svc.GetTasks(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("GetTasks failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestTaskService_GetTaskByID_FromCache(t *testing.T) {
	cachedTask := &models.Task{ID: "t1", UserID: "user-1", Title: "Cached"}
	repo := &mockTaskRepository{
		getByIDFn: func(ctx context.Context, id string) (*models.Task, error) {
			t.Fatal("should not hit DB when cache hit")
			return nil, nil
		},
	}
	cache := &mockTaskCacheRepository{
		getFn: func(ctx context.Context, key string) (*models.Task, error) {
			return cachedTask, nil
		},
	}
	svc := NewTaskService(repo, cache, "app", 10*time.Minute)

	task, err := svc.GetTaskByID(context.Background(), "t1", "user-1")
	if err != nil {
		t.Fatalf("GetTaskByID failed: %v", err)
	}
	if task.Title != "Cached" {
		t.Errorf("expected Cached, got %s", task.Title)
	}
}

func TestTaskService_GetTaskByID_CacheMiss(t *testing.T) {
	dbTask := &models.Task{ID: "t1", UserID: "user-1", Title: "FromDB"}
	setCalled := false
	repo := &mockTaskRepository{
		getByIDFn: func(ctx context.Context, id string) (*models.Task, error) {
			return dbTask, nil
		},
	}
	cache := &mockTaskCacheRepository{
		getFn: func(ctx context.Context, key string) (*models.Task, error) {
			return nil, nil
		},
		setFn: func(ctx context.Context, task *models.Task, key string, exp time.Duration) error {
			setCalled = true
			return nil
		},
	}
	svc := NewTaskService(repo, cache, "app", 10*time.Minute)

	task, err := svc.GetTaskByID(context.Background(), "t1", "user-1")
	if err != nil {
		t.Fatalf("GetTaskByID failed: %v", err)
	}
	if task.Title != "FromDB" {
		t.Errorf("expected FromDB, got %s", task.Title)
	}
	if !setCalled {
		t.Error("expected cache Set to be called on miss")
	}
}

func TestTaskService_GetTaskByID_NotOwner(t *testing.T) {
	repo := &mockTaskRepository{
		getByIDFn: func(ctx context.Context, id string) (*models.Task, error) {
			return &models.Task{ID: "t1", UserID: "other-user", Title: "Secret"}, nil
		},
	}
	cache := &mockTaskCacheRepository{
		getFn: func(ctx context.Context, key string) (*models.Task, error) {
			return nil, nil
		},
	}
	svc := NewTaskService(repo, cache, "app", 10*time.Minute)

	_, err := svc.GetTaskByID(context.Background(), "t1", "user-1")
	if err == nil {
		t.Fatal("expected forbidden error")
	}
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Code != "FORBIDDEN" {
		t.Errorf("expected FORBIDDEN, got %v", err)
	}
}

func TestTaskService_UpdateTaskByID_Success(t *testing.T) {
	deleteCalled := false
	repo := &mockTaskRepository{
		getByIDFn: func(ctx context.Context, id string) (*models.Task, error) {
			return &models.Task{ID: "t1", UserID: "user-1"}, nil
		},
		updateFn: func(ctx context.Context, id string, task *models.Task) error {
			return nil
		},
	}
	cache := &mockTaskCacheRepository{
		getFn: func(ctx context.Context, key string) (*models.Task, error) {
			return nil, nil
		},
		deleteFn: func(ctx context.Context, key string) error {
			deleteCalled = true
			return nil
		},
	}
	svc := NewTaskService(repo, cache, "app", 10*time.Minute)

	err := svc.UpdateTaskByID(context.Background(), "t1", "user-1", &models.Task{Title: "Updated"})
	if err != nil {
		t.Fatalf("UpdateTaskByID failed: %v", err)
	}
	if !deleteCalled {
		t.Error("expected cache invalidation")
	}
}

func TestTaskService_DeleteTaskByID_Success(t *testing.T) {
	deleteCalled := false
	repo := &mockTaskRepository{
		getByIDFn: func(ctx context.Context, id string) (*models.Task, error) {
			return &models.Task{ID: "t1", UserID: "user-1"}, nil
		},
		deleteFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
	cache := &mockTaskCacheRepository{
		getFn: func(ctx context.Context, key string) (*models.Task, error) {
			return nil, nil
		},
		deleteFn: func(ctx context.Context, key string) error {
			deleteCalled = true
			return nil
		},
	}
	svc := NewTaskService(repo, cache, "app", 10*time.Minute)

	err := svc.DeleteTaskByID(context.Background(), "t1", "user-1")
	if err != nil {
		t.Fatalf("DeleteTaskByID failed: %v", err)
	}
	if !deleteCalled {
		t.Error("expected cache deletion")
	}
}

func TestTaskService_DeleteTaskByID_NotFound(t *testing.T) {
	repo := &mockTaskRepository{
		getByIDFn: func(ctx context.Context, id string) (*models.Task, error) {
			return nil, apperror.NewNotFoundError("task not found")
		},
	}
	cache := &mockTaskCacheRepository{
		getFn: func(ctx context.Context, key string) (*models.Task, error) {
			return nil, nil
		},
	}
	svc := NewTaskService(repo, cache, "app", 10*time.Minute)

	err := svc.DeleteTaskByID(context.Background(), "missing", "user-1")
	if err == nil {
		t.Fatal("expected not found error")
	}
}

var _ ports.TaskRepository = (*mockTaskRepository)(nil)
var _ ports.TaskCacheRepository = (*mockTaskCacheRepository)(nil)
