package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	GRPCPort   string `mapstructure:"GRPC_PORT"`

	DBPort     string `mapstructure:"DB_PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisDB       string `mapstructure:"REDIS_DB"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	SessionExpiration time.Duration `mapstructure:"SESSION_EXPIRATION"`
	RedisAppName      string        `mapstructure:"REDIS_APP_NAME"`
	CacheExpiration   time.Duration `mapstructure:"CACHE_EXPIRATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Explicitly bind so Unmarshal sees Docker/K8s environment variables
	for _, key := range []string{
		"SERVER_HOST", "SERVER_PORT", "GRPC_PORT",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_NAME", "DB_PASSWORD",
		"REDIS_ADDR", "REDIS_DB", "REDIS_PASSWORD", "REDIS_APP_NAME",
		"SESSION_EXPIRATION", "CACHE_EXPIRATION",
	} {
		_ = viper.BindEnv(key)
	}

	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")
	viper.SetDefault("GRPC_PORT", "50051")
	viper.SetDefault("REDIS_DB", "0")
	viper.SetDefault("SESSION_EXPIRATION", "30m")
	viper.SetDefault("CACHE_EXPIRATION", "10m")
	viper.SetDefault("REDIS_APP_NAME", "task-management-api")

	// Optional .env file (ignore if missing)
	_ = viper.ReadInConfig()

	err = viper.Unmarshal(&config)
	return config, err
}
