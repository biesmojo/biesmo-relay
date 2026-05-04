package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() error {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil // no DB mode
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return err
	}

	config.MaxConns = 20

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	// Test ping
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return err
	}

	Pool = pool
	return nil
}
