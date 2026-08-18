package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// server
	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	GRPCPort   string `mapstructure:"GRPC_PORT"`

	// postgres
	DBPort     string `mapstructure:"DB_PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	// redis
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisDB       string `mapstructure:"REDIS_DB"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// session
	SessionExpiration time.Duration `mapstructure:"SESSION_EXPIRATION"`
	RedisAppName      string        `mapstructure:"REDIS_APP_NAME"`
	CacheExpiration   time.Duration `mapstructure:"CACHE_EXPIRATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Defaults (used when running under docker-compose env or missing keys)
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")
	viper.SetDefault("GRPC_PORT", "50051")
	viper.SetDefault("REDIS_DB", "0")
	viper.SetDefault("SESSION_EXPIRATION", "30m")
	viper.SetDefault("CACHE_EXPIRATION", "10m")
	viper.SetDefault("REDIS_APP_NAME", "task-management-api")

	// Prefer file if present; ignore if missing (docker-compose injects env)
	_ = viper.ReadInConfig()

	err = viper.Unmarshal(&config)
	return config, err
}
