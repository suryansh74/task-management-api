package service

import (
	"context"
	"errors"
	"testing"

	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
	"github.com/suryansh74/task-management-api-project/internal/utils"
)

func init() {
	logger.Init()
}

// mockUserRepository implements ports.UserRepository
type mockUserRepository struct {
	createUserFn func(ctx context.Context, user *models.User) (string, error)
	findByEmailFn func(ctx context.Context, email string) (*models.User, error)
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, user)
	}
	return "", errors.New("not implemented")
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, errors.New("not implemented")
}

func TestUserService_Register_Success(t *testing.T) {
	repo := &mockUserRepository{
		createUserFn: func(ctx context.Context, user *models.User) (string, error) {
			if user.Name != "Alice" || user.Email != "alice@example.com" {
				t.Errorf("unexpected user data: %+v", user)
			}
			if user.Password == "password123" {
				t.Error("password should be hashed")
			}
			return "user-id-1", nil
		},
	}
	svc := NewUserService(repo)

	resp, err := svc.Register(context.Background(), "Alice", "alice@example.com", "password123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if resp.ID != "user-id-1" || resp.Name != "Alice" || resp.Email != "alice@example.com" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestUserService_Register_RepoError(t *testing.T) {
	repo := &mockUserRepository{
		createUserFn: func(ctx context.Context, user *models.User) (string, error) {
			return "", apperror.NewConflictError("email already exists")
		},
	}
	svc := NewUserService(repo)

	_, err := svc.Register(context.Background(), "Alice", "alice@example.com", "password123")
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Code != "CONFLICT" {
		t.Errorf("expected conflict error, got %v", err)
	}
}

func TestUserService_Login_Success(t *testing.T) {
	hashed, _ := utils.HashedPassword("password123")
	repo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{
				ID:       "user-id-1",
				Name:     "Alice",
				Email:    "alice@example.com",
				Password: hashed,
			}, nil
		},
	}
	svc := NewUserService(repo)

	resp, err := svc.Login(context.Background(), "alice@example.com", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if resp.ID != "user-id-1" || resp.Email != "alice@example.com" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	hashed, _ := utils.HashedPassword("password123")
	repo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{
				ID:       "user-id-1",
				Name:     "Alice",
				Email:    "alice@example.com",
				Password: hashed,
			}, nil
		},
	}
	svc := NewUserService(repo)

	_, err := svc.Login(context.Background(), "alice@example.com", "wrongpass")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Code != "UNAUTHORIZED" {
		t.Errorf("expected unauthorized, got %v", err)
	}
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
			return nil, apperror.NewNotFoundError("user not found")
		},
	}
	svc := NewUserService(repo)

	_, err := svc.Login(context.Background(), "nobody@example.com", "password123")
	if err == nil {
		t.Fatal("expected error")
	}
}

// Ensure interface compliance
var _ ports.UserRepository = (*mockUserRepository)(nil)
