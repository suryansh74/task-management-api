package service

import (
	"context"
	"fmt"
	"time"

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
	return &sessionService{
		sessionRepo:       sessionRepo,
		sessionExpiration: sessionExpiration,
		redisAppName:      redisAppName,
	}
}

// CreateSession it sets new session
// =========================================================================
func (s *sessionService) CreateSession(ctx context.Context, userID string) (string, error) {
	// create random id
	id := utils.MustRandomID()
	sessionID := fmt.Sprintf("%s:sessions:%s", s.redisAppName, id)
	err := s.sessionRepo.Create(ctx, &models.Session{ID: sessionID, UserID: userID}, s.sessionExpiration)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

// Delete it unlink user session
// =========================================================================
func (s *sessionService) Logout(ctx context.Context, sessionID string) error {
	err := s.sessionRepo.Delete(ctx, sessionID)
	if err != nil {
		return err
	}
	return nil
}
