package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

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
	defer conn.Close(ctx)

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
