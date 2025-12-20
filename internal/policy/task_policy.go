package policy

import (
	"context"

	"github.com/suryansh74/task-management-api-project/internal/apperror"
)

type TaskPolicy struct {
	query TaskQuery
}

func NewTaskPolicy(repo TaskQuery) *TaskPolicy {
	return &TaskPolicy{query: repo}
}

func (p *TaskPolicy) mustBeOwner(
	ctx context.Context,
	userID, taskID string,
	key string,
) error {
	if userID == "" {
		return apperror.NewUnauthorizedError("not authenticated")
	}

	TaskUserID, err := p.query.GetOwnerIDByTaskID(ctx, taskID, key)
	if err != nil {
		return err // not found bubbles up
	}

	if TaskUserID != userID {
		return apperror.NewForbiddenError("not allowed")
	}

	return nil
}

// TODO: fix list

func (p *TaskPolicy) CanCreate(ctx context.Context, userID string) error {
	return nil
}

func (p *TaskPolicy) CanRead(ctx context.Context, userID, taskID string, key string) error {
	return p.mustBeOwner(ctx, userID, taskID, key)
}

func (p *TaskPolicy) CanUpdate(ctx context.Context, userID, taskID string, key string) error {
	return p.mustBeOwner(ctx, userID, taskID, key)
}

func (p *TaskPolicy) CanDelete(ctx context.Context, userID, taskID string, key string) error {
	return p.mustBeOwner(ctx, userID, taskID, key)
}
