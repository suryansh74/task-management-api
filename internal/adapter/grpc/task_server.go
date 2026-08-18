package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	taskv1 "github.com/suryansh74/task-management-api-project/api/gen/task/v1"
	"github.com/suryansh74/task-management-api-project/internal/apperror"
	"github.com/suryansh74/task-management-api-project/internal/models"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

// TaskServer is the gRPC adapter for task operations.
// It depends only on the application port (ports.TaskService).
type TaskServer struct {
	taskv1.UnimplementedTaskServiceServer
	taskService ports.TaskService
}

// NewTaskServer creates a new gRPC task server adapter.
func NewTaskServer(taskService ports.TaskService) *TaskServer {
	return &TaskServer{taskService: taskService}
}

func (s *TaskServer) GetTasks(ctx context.Context, req *taskv1.GetTasksRequest) (*taskv1.GetTasksResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	tasks, err := s.taskService.GetTasks(ctx, req.UserId)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &taskv1.GetTasksResponse{
		Tasks: make([]*taskv1.Task, 0, len(tasks)),
	}
	for _, t := range tasks {
		resp.Tasks = append(resp.Tasks, toProtoTask(t))
	}
	return resp, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.GetTaskResponse, error) {
	if req.Id == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "id and user_id are required")
	}

	task, err := s.taskService.GetTaskByID(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, mapError(err)
	}

	return &taskv1.GetTaskResponse{Task: toProtoTask(task)}, nil
}

func (s *TaskServer) CreateTask(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.CreateTaskResponse, error) {
	if req.UserId == "" || req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and title are required")
	}

	task := &models.Task{
		UserID:  req.UserId,
		Title:   req.Title,
		Content: req.Content,
	}

	id, err := s.taskService.CreateTask(ctx, task)
	if err != nil {
		return nil, mapError(err)
	}

	// Re-fetch so we return the full created entity (timestamps, etc.)
	created, err := s.taskService.GetTaskByID(ctx, id, req.UserId)
	if err != nil {
		// Creation succeeded; return at least the id
		return &taskv1.CreateTaskResponse{Id: id}, nil
	}

	return &taskv1.CreateTaskResponse{
		Id:   id,
		Task: toProtoTask(created),
	}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.UpdateTaskResponse, error) {
	if req.Id == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "id and user_id are required")
	}

	task := &models.Task{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := s.taskService.UpdateTaskByID(ctx, req.Id, req.UserId, task); err != nil {
		return nil, mapError(err)
	}

	updated, err := s.taskService.GetTaskByID(ctx, req.Id, req.UserId)
	if err != nil {
		return &taskv1.UpdateTaskResponse{}, nil
	}

	return &taskv1.UpdateTaskResponse{Task: toProtoTask(updated)}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *taskv1.DeleteTaskRequest) (*taskv1.DeleteTaskResponse, error) {
	if req.Id == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "id and user_id are required")
	}

	if err := s.taskService.DeleteTaskByID(ctx, req.Id, req.UserId); err != nil {
		return nil, mapError(err)
	}

	return &taskv1.DeleteTaskResponse{}, nil
}

func toProtoTask(t *models.Task) *taskv1.Task {
	if t == nil {
		return nil
	}
	return &taskv1.Task{
		Id:        t.ID,
		UserId:    t.UserID,
		Title:     t.Title,
		Content:   t.Content,
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
	}
}

// mapError converts application errors into gRPC status errors.
func mapError(err error) error {
	if err == nil {
		return nil
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case "NOT_FOUND":
			return status.Error(codes.NotFound, appErr.Message)
		case "UNAUTHORIZED", "FORBIDDEN":
			return status.Error(codes.PermissionDenied, appErr.Message)
		case "VALIDATION_ERROR", "BAD_REQUEST":
			return status.Error(codes.InvalidArgument, appErr.Message)
		case "CONFLICT":
			return status.Error(codes.AlreadyExists, appErr.Message)
		default:
			return status.Error(codes.Internal, appErr.Message)
		}
	}

	return status.Error(codes.Internal, err.Error())
}
