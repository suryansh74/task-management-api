package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) ports.UserRepository {
	logger.Log.Info().Msg("initializing user repository")
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	logger.Log.Debug().
		Str("email", user.Email).
		Str("name", user.Name).
		Msg("creating new user")

	var id string
	err := ur.db.QueryRow(ctx,
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			logger.Log.Warn().
				Str("email", user.Email).
				Str("pg_error_code", pgErr.Code).
				Msg("duplicate user registration attempt")
			return "", apperror.NewConflictError("User with this email already exists")
		}
		logger.Log.Error().
			Err(err).
			Str("email", user.Email).
			Str("name", user.Name).
			Msg("failed to create user")
		return "", apperror.NewInternalError("Failed to create user", err)
	}

	logger.Log.Info().
		Str("user_id", id).
		Str("email", user.Email).
		Str("name", user.Name).
		Msg("user created successfully")
	return id, nil
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	logger.Log.Debug().
		Str("email", email).
		Msg("finding user by email")

	var user models.User
	err := ur.db.QueryRow(ctx,
		"SELECT id, name, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Log.Warn().
				Str("email", email).
				Msg("user not found")
			return nil, apperror.NewNotFoundError("User not found")
		}
		logger.Log.Error().
			Err(err).
			Str("email", email).
			Msg("failed to find user by email")
		return nil, apperror.NewInternalError("Failed to retrieve user", err)
	}

	logger.Log.Info().
		Str("user_id", user.ID).
		Str("email", email).
		Msg("user found successfully")
	return &user, nil
}
