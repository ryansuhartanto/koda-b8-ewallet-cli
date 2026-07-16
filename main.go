package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file", err)
	}

	_, err = pgx.Connect(context.Background(), "")
	if err != nil {
		log.Panic("Error connecting to database", err)
	}
}
