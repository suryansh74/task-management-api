package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type sessionRepository struct {
	redisClient *redis.Client
}

// NewSessionRepository constructor for session repository
// =========================================================================
func NewSessionRepository(redisClient *redis.Client) ports.SessionRepository {
	return &sessionRepository{
		redisClient: redisClient,
	}
}

// Create it set session
// =========================================================================
func (us *sessionRepository) Create(ctx context.Context, session *models.Session, sessionExpiration time.Duration) error {
	err := us.redisClient.HSet(ctx, session.ID, models.Session{ID: session.ID, UserID: session.UserID}).Err()
	if err != nil {
		logger.Log.Err(err)
		return apperror.NewInternalError("unable to set user session in redis", err)
	}
	err = us.redisClient.Expire(ctx, session.ID, sessionExpiration).Err()
	if err != nil {
		logger.Log.Err(err)
		return apperror.NewInternalError("unable to set expiry for user session in redis", err)
	}
	logger.Log.Info().Str("session_id", session.ID).Msg("user session created")
	return nil
}

// GetByID finds session
// =========================================================================
func (us *sessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	return nil, nil
}

// Delete it unlink session
// =========================================================================
func (us *sessionRepository) Delete(ctx context.Context, sessionID string) error {
	err := us.redisClient.Unlink(ctx, sessionID).Err()
	if err != nil {
		logger.Log.Err(err).Msg("failed to unlink session")
		return apperror.NewInternalError("unable to delete session", err)
	}
	return nil
}
