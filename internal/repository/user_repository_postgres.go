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
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	var id string
	err := ur.db.QueryRow(ctx,
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			logger.Log.Warn().Str("email", user.Email).Msg("Duplicate user registration attempt")
			return "", apperror.NewConflictError("User with this email already exists")
		}
		logger.Log.Error().Err(err).Msg("Failed to create user")
		return "", apperror.NewInternalError("Failed to create user", err)
	}

	return id, nil
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := ur.db.QueryRow(ctx,
		"SELECT id, name, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFoundError("User not found")
		}
		logger.Log.Error().Err(err).Msg("Failed to find user")
		return nil, apperror.NewInternalError("Failed to retrieve user", err)
	}

	return &user, nil
}
