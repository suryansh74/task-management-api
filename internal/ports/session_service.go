package ports

import (
	"context"
)

// SessionService defines business logic operations for users
type SessionService interface {
	CreateSession(ctx context.Context, userID string) (string, error)
	Logout(ctx context.Context, sessionID string) error
}
