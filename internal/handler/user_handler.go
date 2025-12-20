package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/http/response"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/ports"
	"github.com/suryansh74/task-management-api-project/internal/validator"
)

type UserHandler struct {
	userService       ports.UserService
	sessionService    ports.SessionService
	sessionExpiration time.Duration
	redisAppName      string
}

// NewUserHandler Constructor for UserHandler
// =========================================================================
func NewUserHandler(userService ports.UserService, sessionService ports.SessionService, sessionExpiration time.Duration, redisAppName string) *UserHandler {
	logger.Log.Info().
		Dur("session_expiration", sessionExpiration).
		Str("redis_app_name", redisAppName).
		Msg("initializing user handler")
	return &UserHandler{
		userService:       userService,
		sessionService:    sessionService,
		sessionExpiration: sessionExpiration,
		redisAppName:      redisAppName,
	}
}

// RegisterRequest dto for incoming req
// =========================================================================
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// LoginRequest dto for incoming req
// =========================================================================
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Register handles user registration
// =========================================================================
func (h *UserHandler) Register(c *fiber.Ctx) error {
	logger.Log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("ip", c.IP()).
		Msg("received user registration request")

	var req RegisterRequest

	// Parse body
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Warn().
			Err(err).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Msg("failed to parse registration request body")
		return apperror.NewBadRequestError("Invalid request body")
	}

	logger.Log.Debug().
		Str("email", req.Email).
		Str("name", req.Name).
		Msg("parsed registration request")

	// Validate
	if fieldErrors := validator.ValidateStruct(req); len(fieldErrors) > 0 {
		logger.Log.Warn().
			Interface("validation_errors", fieldErrors).
			Str("email", req.Email).
			Str("name", req.Name).
			Msg("validation failed for registration")
		return response.ValidationError(c, fieldErrors)
	}

	logger.Log.Debug().
		Str("email", req.Email).
		Str("name", req.Name).
		Msg("registering new user")

	// Call service
	user, err := h.userService.Register(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("email", req.Email).
			Str("name", req.Name).
			Msg("failed to register user")
		return err // Global error handler will catch this
	}

	logger.Log.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("name", user.Name).
		Int("status", fiber.StatusCreated).
		Msg("user registered successfully")

	return response.Success(c, fiber.StatusCreated, "User registered successfully.", user)
}

// Login handles retrieving user by email
// =========================================================================
func (h *UserHandler) Login(c *fiber.Ctx) error {
	logger.Log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("ip", c.IP()).
		Msg("received user login request")

	var req LoginRequest

	// Parse body
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Warn().
			Err(err).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Msg("failed to parse login request body")
		return apperror.NewBadRequestError("Invalid request body")
	}

	logger.Log.Debug().
		Str("email", req.Email).
		Msg("parsed login request")

	// Validate
	if fieldErrors := validator.ValidateStruct(req); len(fieldErrors) > 0 {
		logger.Log.Warn().
			Interface("validation_errors", fieldErrors).
			Str("email", req.Email).
			Msg("validation failed for login")
		return response.ValidationError(c, fieldErrors)
	}

	logger.Log.Debug().
		Str("email", req.Email).
		Msg("attempting user login")

	// Call service
	user, err := h.userService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		logger.Log.Warn().
			Err(err).
			Str("email", req.Email).
			Str("ip", c.IP()).
			Msg("login attempt failed")
		return err
	}

	logger.Log.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Msg("user authenticated successfully, creating session")

	// set user session
	sessionID, err := h.sessionService.CreateSession(c.Context(), user.ID)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("user_id", user.ID).
			Str("email", user.Email).
			Msg("failed to create user session")
		return apperror.NewInternalError("user session not created", err)
	}

	logger.Log.Debug().
		Str("user_id", user.ID).
		Str("session_id", sessionID).
		Dur("expiration", h.sessionExpiration).
		Msg("session created, setting cookie")

	// set session id in HTTP-ONLY cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(h.sessionExpiration),
	})

	logger.Log.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("session_id", sessionID).
		Int("status", fiber.StatusOK).
		Msg("user logged in successfully with session")

	return response.Success(c, fiber.StatusOK, "User logged in successfully", user)
}

// Logout delete session
// =========================================================================
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")

	logger.Log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("session_id", sessionID).
		Str("ip", c.IP()).
		Msg("received user logout request")

	if sessionID == "" {
		logger.Log.Debug().
			Str("ip", c.IP()).
			Msg("logout request with no session cookie")
	} else {
		logger.Log.Debug().
			Str("session_id", sessionID).
			Msg("deleting user session")

		// Call service
		err := h.sessionService.Logout(c.Context(), sessionID)
		if err != nil {
			logger.Log.Warn().
				Err(err).
				Str("session_id", sessionID).
				Msg("failed to delete session, continuing with logout")
		} else {
			logger.Log.Info().
				Str("session_id", sessionID).
				Msg("session deleted successfully")
		}
	}

	logger.Log.Debug().
		Str("session_id", sessionID).
		Msg("clearing session cookie")

	// set session id in HTTP-ONLY cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	logger.Log.Info().
		Str("session_id", sessionID).
		Int("status", fiber.StatusOK).
		Msg("user logged out successfully")

	return response.Success(c, fiber.StatusOK, "User logout successfully", nil)
}
