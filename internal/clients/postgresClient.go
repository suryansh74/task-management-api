package clients

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func PostgresClient(user, password, host, port, dbName string) *pgx.Conn {
	dbPath := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbName,
	)

	var conn *pgx.Conn
	var err error
	for attempt := 1; attempt <= 20; attempt++ {
		conn, err = pgx.Connect(context.Background(), dbPath)
		if err == nil {
			fmt.Fprintf(os.Stderr, "PostgreSQL connected (attempt %d)\n", attempt)
			return conn
		}
		fmt.Fprintf(os.Stderr, "PostgreSQL dial attempt %d/20 failed: %v (retry in 2s)\n", attempt, err)
		time.Sleep(2 * time.Second)
	}
	fmt.Fprintf(os.Stderr, "Unable to connect to database after retries: %v\n", err)
	os.Exit(1)
	return nil
}
