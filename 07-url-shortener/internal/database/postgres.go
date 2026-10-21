package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create PostgreSQL pool: %w", err,
		)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"failed to connect to PostgreSQL: %w", err,
		)
	}

	return pool, nil
}
