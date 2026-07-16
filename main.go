package main

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

//go:embed db/migrations/*
var migrationsFS embed.FS

func init() {
	d, err := iofs.New(migrationsFS, "db/migrations")
	if err != nil {
		log.Panicln("Error creating FS driver", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, "pgx5://:")
	if err != nil {
		log.Panicln("Error creating migration instance", err)
	}

	err = m.Up()
	if err != nil {
		log.Panicln("Error migrating database", err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Error loading .env file", err)
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "")
	if err != nil {
		log.Panicln("Error connecting to database", err)
	}

	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' ORDER BY table_name;`
	rows, err := conn.Query(ctx, query)
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
