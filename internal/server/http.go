package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	grpcadapter "github.com/suryansh74/task-management-api-project/internal/adapter/grpc"
	"github.com/suryansh74/task-management-api-project/internal/config"
	"github.com/suryansh74/task-management-api-project/internal/handler"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/ports"
	"github.com/suryansh74/task-management-api-project/internal/repository"
	"github.com/suryansh74/task-management-api-project/internal/service"
)

type server struct {
	app            *fiber.App
	redisClient    *redis.Client
	postgresClient *pgx.Conn
	cfg            *config.Config
}

// StartServer wires repositories → services → adapters (REST + gRPC) and starts both servers.
func StartServer(app *fiber.App, redisClient *redis.Client, postgresClient *pgx.Conn, cfg *config.Config) {
	server := &server{
		app:            app,
		redisClient:    redisClient,
		postgresClient: postgresClient,
		cfg:            cfg,
	}

	// Initialize repositories (driven adapters)
	var userRepo ports.UserRepository = repository.NewUserRepository(postgresClient)
	var sessionRepo ports.SessionRepository = repository.NewSessionRepository(redisClient)
	var taskRepo ports.TaskRepository = repository.NewTaskRepository(postgresClient)
	var taskCacheRepo ports.TaskCacheRepository = repository.NewTaskCacheRepository(redisClient)

	// Initialize services (application core)
	var userService ports.UserService = service.NewUserService(userRepo)
	var sessionService ports.SessionService = service.NewSessionService(sessionRepo, cfg.SessionExpiration, cfg.RedisAppName)
	var taskService ports.TaskService = service.NewTaskService(taskRepo, taskCacheRepo, cfg.RedisAppName, cfg.CacheExpiration)

	// Initialize HTTP handlers (driving adapters – REST)
	var userHandler ports.UserHandler = handler.NewUserHandler(userService, sessionService, cfg.SessionExpiration, cfg.RedisAppName)
	var taskHandler ports.TaskHandler = handler.NewTaskHandler(taskService, cfg.RedisAppName, cfg.SessionExpiration)

	server.setupRoutes(userHandler, taskHandler)

	// Initialize gRPC server (driving adapter – gRPC)
	// Shares the same taskService instance as REST
	grpcPort := cfg.GRPCPort
	if grpcPort == "" {
		grpcPort = "50051"
	}
	grpcAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, grpcPort)
	grpcServer := grpcadapter.NewServer(grpcAddr, taskService)

	// Start gRPC in background
	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Start REST (blocks)
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	logger.Log.Info().Msg("REST server starting on " + addr)
	if err := app.Listen(addr); err != nil {
		logger.Log.Fatal().Err(err).Msg("REST server failed")
	}
}
