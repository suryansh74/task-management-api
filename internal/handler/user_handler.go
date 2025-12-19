package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/http/response"
	"github.com/suryansh74/task-management-api-project/internal/ports"
	"github.com/suryansh74/task-management-api-project/internal/validator"
)

type UserHandler struct {
	userService ports.UserService
}

// NewUserHandler Constructor for UserHandler
// =========================================================================
func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest dto for incoming req
// =========================================================================
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type GetUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Register handles user registration
// =========================================================================
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req CreateUserRequest

	// Parse body
	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequestError("Invalid request body")
	}

	// Validate
	if fieldErrors := validator.ValidateStruct(req); len(fieldErrors) > 0 {
		return response.ValidationError(c, fieldErrors)
	}

	// Call service
	user, err := h.userService.Register(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return err // Global error handler will catch this
	}

	return response.Success(c, fiber.StatusCreated, "User registered successfully.", user)
}

// Login handles retrieving user by email
// =========================================================================
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req GetUserRequest

	// Parse body
	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequestError("Invalid request body")
	}

	// Validate
	if fieldErrors := validator.ValidateStruct(req); len(fieldErrors) > 0 {
		return response.ValidationError(c, fieldErrors)
	}

	// Call service
	user, err := h.userService.Login(c.Context(), req.Email)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "User logged in successfully", user)
}
