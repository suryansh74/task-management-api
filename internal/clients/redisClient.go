package clients

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisClient(addr string, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        password,
		DB:              db,
		PoolSize:        10,
		MinIdleConns:    2,
		ConnMaxIdleTime: 5 * time.Second,
		DialTimeout:     5 * time.Second,
		WriteTimeout:    3 * time.Second,
		ReadTimeout:     3 * time.Second,
		MaxRetries:      3,
	})

	var err error
	for attempt := 1; attempt <= 20; attempt++ {
		err = rdb.Ping(context.Background()).Err()
		if err == nil {
			fmt.Fprintf(os.Stderr, "Redis connected (attempt %d)\n", attempt)
			return rdb
		}
		fmt.Fprintf(os.Stderr, "Redis ping attempt %d/20 failed: %v (retry in 2s)\n", attempt, err)
		time.Sleep(2 * time.Second)
	}
	fmt.Fprintf(os.Stderr, "Unable to connect to Redis after retries: %v\n", err)
	os.Exit(1)
	return nil
}
