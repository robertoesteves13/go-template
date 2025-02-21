package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	conn *pgxpool.Pool
)

// Initializes the database connection pool
func ConnectDatabase() error {
	var err error
	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		return fmt.Errorf("env variable `DATABASE_URL` not set")
	}

	conn, err = pgxpool.New(context.Background(), db_url)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	return nil
}

// Get a connection from the pool
func GetConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return conn.Acquire(ctx)
}

// Closes the pool. All connections will be closed aswell.
func CloseConn() {
	conn.Close()
}
