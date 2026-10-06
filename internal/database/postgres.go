package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDatabase(databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(
		context.Background(),
		databaseURL,
	)

	if err != nil {
		return pool, nil
	}

	return pool, nil
}
