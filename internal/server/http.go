package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/suryansh74/task-management-api-project/internal/config"
	"github.com/suryansh74/task-management-api-project/internal/handler"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/repository"
	"github.com/suryansh74/task-management-api-project/internal/service"
)

type server struct {
	app            *fiber.App
	redisClient    *redis.Client
	postgresClient *pgx.Conn
}

func StartServer(app *fiber.App, redisClient *redis.Client, postgresClient *pgx.Conn, cfg config.Config) {
	server := &server{
		app:            app,
		redisClient:    redisClient,
		postgresClient: postgresClient,
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(postgresClient)

	// Initialize services (no config needed now)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	server.setupRoutes(userHandler)
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	logger.Log.Info().Msg("server starting on port:" + addr)
	app.Listen(addr)
}
