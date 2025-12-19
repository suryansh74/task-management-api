package clients

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisClient(addr string, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		// setting for connecting to rds
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB

		// pool settng
		PoolSize:        10,
		MinIdleConns:    2,
		ConnMaxIdleTime: 5 * time.Second, // how long connection will wait unused

		// timeout
		DialTimeout:  5 * time.Second, // time duration for connecting to database and creating new instance
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,

		// retris
		MaxRetries: 3,
	})

	return rdb
}
