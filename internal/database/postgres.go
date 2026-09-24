package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDatabase(databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(
		context.Background(),
		databaseURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	return pool, nil
}
