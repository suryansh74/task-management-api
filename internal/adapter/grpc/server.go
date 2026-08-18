package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	taskv1 "github.com/suryansh74/task-management-api-project/api/gen/task/v1"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/ports"
)

// Server wraps the gRPC server.
type Server struct {
	grpcServer *grpc.Server
	addr       string
}

// NewServer creates a gRPC server with the Task service registered.
// Both REST and gRPC share the same ports.TaskService instance.
func NewServer(addr string, taskService ports.TaskService) *Server {
	s := grpc.NewServer()

	taskServer := NewTaskServer(taskService)
	taskv1.RegisterTaskServiceServer(s, taskServer)

	// Register reflection for tools like grpcurl
	reflection.Register(s)

	return &Server{
		grpcServer: s,
		addr:       addr,
	}
}

// Start begins serving gRPC on the configured address.
// This call blocks; run it in a goroutine from main.
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}

	logger.Log.Info().Msg("gRPC server starting on " + s.addr)
	return s.grpcServer.Serve(lis)
}

// Stop gracefully stops the gRPC server.
func (s *Server) Stop() {
	logger.Log.Info().Msg("stopping gRPC server")
	s.grpcServer.GracefulStop()
}
