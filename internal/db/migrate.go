package db

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*
var migrationsFS embed.FS

func Migrate(pool *pgxpool.Pool) (*migrate.Migrate, error) {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("migrate source: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		source.Close()
		db.Close()
		return nil, fmt.Errorf("migrate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "pgx5", driver)
	if err != nil {
		driver.Close()
		source.Close()
		return nil, fmt.Errorf("migrate instance: %w", err)
	}

	return m, nil
}
