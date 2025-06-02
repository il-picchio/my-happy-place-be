package db

import (
	"blenioviva/internal/db/generated"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

func New() *DB {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		panic("DATABASE_URL not set")
	}
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		panic(fmt.Sprintf("unable to connect to database: %v", err))
	}
	return &DB{
		pool:    pool,
		queries: generated.New(pool),
	}
}
