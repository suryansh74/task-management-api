package ports

import (
	"context"
	"time"

	"github.com/suryansh74/task-management-api-project/internal/models"
)

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session, sessionExpiration time.Duration) error
	GetByID(ctx context.Context, id string) (*models.Session, error)
	Delete(ctx context.Context, sessionID string) error
}
