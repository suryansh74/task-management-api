package apperror

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("not logged in")
	if err.Code != "UNAUTHORIZED" {
		t.Errorf("expected code UNAUTHORIZED, got %s", err.Code)
	}
	if err.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, err.StatusCode)
	}
	if err.Message != "not logged in" {
		t.Errorf("unexpected message: %s", err.Message)
	}
	if !errors.Is(err, ErrUnauthorized) {
		t.Error("expected errors.Is to match ErrUnauthorized")
	}
}

func TestNewForbiddenError(t *testing.T) {
	err := NewForbiddenError("not allowed")
	if err.Code != "FORBIDDEN" || err.StatusCode != http.StatusForbidden {
		t.Errorf("unexpected forbidden error: %+v", err)
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("task not found")
	if err.Code != "NOT_FOUND" || err.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected not found error: %+v", err)
	}
}

func TestNewConflictError(t *testing.T) {
	err := NewConflictError("email exists")
	if err.Code != "CONFLICT" || err.StatusCode != http.StatusConflict {
		t.Errorf("unexpected conflict error: %+v", err)
	}
}

func TestNewBadRequestError(t *testing.T) {
	err := NewBadRequestError("bad body")
	if err.Code != "BAD_REQUEST" || err.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected bad request error: %+v", err)
	}
}

func TestNewInternalError(t *testing.T) {
	inner := errors.New("db down")
	err := NewInternalError("something failed", inner)
	if err.Code != "INTERNAL_ERROR" || err.StatusCode != http.StatusInternalServerError {
		t.Errorf("unexpected internal error: %+v", err)
	}
	if !errors.Is(err, inner) {
		t.Error("expected unwrap to return inner error")
	}
}

func TestAppError_Error(t *testing.T) {
	err := NewUnauthorizedError("auth failed")
	if err.Error() == "" {
		t.Error("Error() should not return empty string")
	}
}
