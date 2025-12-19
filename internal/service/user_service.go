package service

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
	"github.com/suryansh74/task-management-api-project/internal/utils"
)

type userService struct {
	userRepo ports.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user account
func (s *userService) CreateUser(ctx context.Context, name, email, password string) (*ports.UserResponse, error) {
	// Hash password
	hashedPassword, err := utils.HashedPassword(password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		return nil, apperror.NewInternalError("Failed to process password", err)
	}

	// Create user model
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	// Save to repository
	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err // Repository already returns AppError
	}

	logger.Log.Info().Str("user_id", userID).Str("email", email).Msg("User created successfully")

	return &ports.UserResponse{
		ID:    userID,
		Name:  name,
		Email: email,
	}, nil
}

// GetUserByEmail retrieves user information by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*ports.UserResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err // Repository already returns AppError
	}

	logger.Log.Info().Str("user_id", user.ID).Str("email", email).Msg("User retrieved successfully")

	return &ports.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
