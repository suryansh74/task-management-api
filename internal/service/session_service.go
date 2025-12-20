package service

import (
	"context"
	"fmt"
	"time"

	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
	"github.com/suryansh74/task-management-api-project/internal/utils"
)

type sessionService struct {
	sessionRepo       ports.SessionRepository
	sessionExpiration time.Duration
	redisAppName      string
}

// NewSessionService creates a new user session service instance
// =========================================================================
func NewSessionService(sessionRepo ports.SessionRepository, sessionExpiration time.Duration, redisAppName string) ports.SessionService {
	logger.Log.Info().
		Dur("session_expiration", sessionExpiration).
		Str("redis_app_name", redisAppName).
		Msg("initializing session service")
	return &sessionService{
		sessionRepo:       sessionRepo,
		sessionExpiration: sessionExpiration,
		redisAppName:      redisAppName,
	}
}

// CreateSession it sets new session
// =========================================================================
func (s *sessionService) CreateSession(ctx context.Context, userID string) (string, error) {
	logger.Log.Debug().
		Str("user_id", userID).
		Msg("creating new session for user")

	// create random id
	id := utils.MustRandomID()
	sessionID := fmt.Sprintf("%s:sessions:%s", s.redisAppName, id)

	logger.Log.Debug().
		Str("user_id", userID).
		Str("session_id", sessionID).
		Str("random_id", id).
		Msg("generated session id")

	err := s.sessionRepo.Create(ctx, &models.Session{ID: sessionID, UserID: userID}, s.sessionExpiration)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("user_id", userID).
			Str("session_id", sessionID).
			Msg("failed to create session")
		return "", err
	}

	logger.Log.Info().
		Str("user_id", userID).
		Str("session_id", sessionID).
		Dur("expiration", s.sessionExpiration).
		Msg("session created successfully")
	return sessionID, nil
}

// Delete it unlink user session
// =========================================================================
func (s *sessionService) Logout(ctx context.Context, sessionID string) error {
	logger.Log.Debug().
		Str("session_id", sessionID).
		Msg("logging out user session")

	err := s.sessionRepo.Delete(ctx, sessionID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("failed to logout user session")
		return err
	}

	logger.Log.Info().
		Str("session_id", sessionID).
		Msg("user logged out successfully")
	return nil
}
