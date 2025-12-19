package ports

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}
