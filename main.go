package main

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/suryansh74/task-management-api-project/internal/clients"
	"github.com/suryansh74/task-management-api-project/internal/config"
	"github.com/suryansh74/task-management-api-project/internal/logger"
	"github.com/suryansh74/task-management-api-project/internal/server"
)

func main() {
	// initialize logger
	logger.Init()
	logger.Log.Info().Msg("Application starting")

	// load config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Cannot load config")
	}

	// postgres setup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	postgresClient := clients.PostgresClient(
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	defer postgresClient.Close(ctx)
	logger.Log.Info().Msg("PostgreSQL connected")

	// redis setup
	redisDB, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Invalid REDIS_DB")
	}

	redisClient := clients.RedisClient(
		cfg.RedisAddr,
		cfg.RedisPassword,
		redisDB,
	)
	defer redisClient.Close()
	logger.Log.Info().Msg("Redis connected")

	// fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorHandler: server.ErrorHandler(),
	})
	app.Use(server.RequestLogger())

	// routes
	server.StartServer(app, redisClient, postgresClient, &cfg)
}
