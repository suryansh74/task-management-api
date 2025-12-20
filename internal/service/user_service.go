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
	logger.Log.Info().Msg("initializing user service")
	return &userService{
		userRepo: userRepo,
	}
}

// Register creates a new user account
// =========================================================================
func (s *userService) Register(ctx context.Context, name, email, password string) (*ports.UserResponse, error) {
	logger.Log.Debug().
		Str("email", email).
		Str("name", name).
		Msg("registering new user")

	// Hash password
	logger.Log.Debug().
		Str("email", email).
		Msg("hashing user password")

	hashedPassword, err := utils.HashedPassword(password)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("email", email).
			Msg("failed to hash password")
		return nil, apperror.NewInternalError("Failed to process password", err)
	}

	logger.Log.Debug().
		Str("email", email).
		Msg("password hashed successfully")

	// Create user model
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	// NOTE: no need to check if user already exists repo check itself while creating new user

	// Save to repository
	logger.Log.Debug().
		Str("email", email).
		Str("name", name).
		Msg("saving user to repository")

	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("email", email).
			Str("name", name).
			Msg("failed to create user in repository")
		return nil, err // Repository already returns AppError
	}

	logger.Log.Info().
		Str("user_id", userID).
		Str("email", email).
		Str("name", name).
		Msg("user registered successfully")

	return &ports.UserResponse{
		ID:    userID,
		Name:  name,
		Email: email,
	}, nil
}

// Login retrieves user information by email
// =========================================================================
func (s *userService) Login(ctx context.Context, email string, password string) (*ports.UserResponse, error) {
	logger.Log.Debug().
		Str("email", email).
		Msg("attempting user login")

	// Find user by email
	logger.Log.Debug().
		Str("email", email).
		Msg("finding user by email")

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.Log.Warn().
			Err(err).
			Str("email", email).
			Msg("user not found or repository error")
		return nil, err // Repository already returns AppError
	}

	logger.Log.Debug().
		Str("user_id", user.ID).
		Str("email", email).
		Msg("user found, verifying password")

	// check hash matches
	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		logger.Log.Warn().
			Str("user_id", user.ID).
			Str("email", email).
			Msg("invalid password attempt")
		return nil, apperror.NewUnauthorizedError("invalid email or password")
	}

	logger.Log.Info().
		Str("user_id", user.ID).
		Str("email", email).
		Msg("user logged in successfully")

	return &ports.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
