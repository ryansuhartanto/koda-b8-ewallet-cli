package main

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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
