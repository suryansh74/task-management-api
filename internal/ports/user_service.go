package ports

import (
	"context"
)

// UserService defines business logic operations for users
type UserService interface {
	Register(ctx context.Context, name, email, password string) (*UserResponse, error)
	Login(ctx context.Context, email string) (*UserResponse, error)
}

// UserResponse is the service layer response for user data
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
