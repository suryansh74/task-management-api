package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/suryansh74/task-management-api-project/internal/config"
	"github.com/suryansh74/task-management-api-project/internal/logger"
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

	server.setupRoutes()
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	logger.Log.Info().Msg("server starting on port:" + addr)
	app.Listen(addr)
}
