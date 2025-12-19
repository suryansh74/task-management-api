package ports

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id string) (*models.Session, error)
	Delete(ctx context.Context, id string) error
}
