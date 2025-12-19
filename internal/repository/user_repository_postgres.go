// Package repository postgres query for users and tasks
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

type userRepository struct {
	db *pgx.Conn
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

// NewUserRepository return postgres user repository
func NewUserRepository(db *pgx.Conn) ports.UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateUser
// ==================================================
func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	var id string
	err := ur.db.QueryRow(ctx, "insert into users (name, email, password) values ($1, $2, $3) returning id", user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			logger.Log.Err(err)
			return "", ErrUserAlreadyExists
		}
		logger.Log.Err(err)
		return "", err
	}
	return id, nil
}

// FindByEmail
// ==================================================
func (ur *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := ur.db.QueryRow(ctx,
		"select id, name, email, password from users where email = $1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Log.Err(err)
			return nil, ErrUserNotFound
		}
		logger.Log.Err(err)
		return nil, err
	}

	return &user, nil
}
