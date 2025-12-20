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
	logger.Log.Info().Msg("initializing session repository")
	return &sessionRepository{
		redisClient: redisClient,
	}
}

// Create it set session
// =========================================================================
func (us *sessionRepository) Create(ctx context.Context, session *models.Session, sessionExpiration time.Duration) error {
	logger.Log.Debug().
		Str("session_id", session.ID).
		Str("user_id", session.UserID).
		Dur("expiration", sessionExpiration).
		Msg("creating user session")

	err := us.redisClient.HSet(ctx, session.ID, models.Session{ID: session.ID, UserID: session.UserID}).Err()
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", session.ID).
			Str("user_id", session.UserID).
			Msg("failed to set user session in redis")
		return apperror.NewInternalError("unable to set user session in redis", err)
	}

	logger.Log.Debug().
		Str("session_id", session.ID).
		Msg("session created successfully, setting expiration")

	err = us.redisClient.Expire(ctx, session.ID, sessionExpiration).Err()
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", session.ID).
			Dur("expiration", sessionExpiration).
			Msg("failed to set expiry for user session in redis")
		return apperror.NewInternalError("unable to set expiry for user session in redis", err)
	}

	logger.Log.Info().
		Str("session_id", session.ID).
		Str("user_id", session.UserID).
		Dur("expiration", sessionExpiration).
		Msg("user session created successfully")
	return nil
}

// GetByID finds session
// =========================================================================
func (us *sessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	logger.Log.Debug().
		Str("session_id", id).
		Msg("retrieving session by id")
	return nil, nil
}

// Delete it unlink session
// =========================================================================
func (us *sessionRepository) Delete(ctx context.Context, sessionID string) error {
	logger.Log.Debug().
		Str("session_id", sessionID).
		Msg("deleting user session")

	err := us.redisClient.Unlink(ctx, sessionID).Err()
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("failed to unlink session from redis")
		return apperror.NewInternalError("unable to delete session", err)
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("user session deleted successfully")
	return nil
}
