package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type sessionRepository struct {
	redisClient *redis.Client
}

func NewSessionRepository(redisClient *redis.Client) ports.SessionRepository {
	return &sessionRepository{
		redisClient: redisClient,
	}
}

func (us *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	return nil
}

func (us *sessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	return nil, nil
}

func (us *sessionRepository) Delete(ctx context.Context, id string) error {
	return nil
}
