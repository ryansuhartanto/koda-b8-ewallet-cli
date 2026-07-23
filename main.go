package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Error loading .env file", err)
	}
}

var ctx context.Context
var pool *pgxpool.Pool

func init() {
	var err error
	ctx = context.Background()
	pool, err = pgxpool.New(ctx, "")
	if err != nil {
		log.Panicln("Error connecting to database", err)
	}

	m, err := db.Migrate(pool)
	if err != nil {
		log.Panicln("Error creating migration", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Panicln("Error migrating database", err)
	}
}

func main() {
	defer pool.Close()

	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' ORDER BY table_name;`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalln("Query failed", err)
	}
	defer rows.Close()

	tables, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		log.Fatalln("Failed to collect rows", err)
	}

	fmt.Println("Tables:", tables)
}
