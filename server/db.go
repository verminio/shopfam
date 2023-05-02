package server

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

func DB(dbUrl string) (*sql.DB, error) {
	if _, err := os.Stat(dbUrl); err != nil {
		if err := createDatabaseFile(dbUrl); err != nil {
			return nil, err
		}
	}

	instance, err := sql.Open("sqlite3", dbUrl)

	if err != nil {
		return nil, fmt.Errorf("failed to open db connection. %w", err)
	}

	return instance, nil
}

func createDatabaseFile(dbUrl string) error {
	dir := path.Dir(dbUrl)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create db directory %s: %w", dir, err)
	}

	if _, err := os.Create(dbUrl); err != nil {
		return fmt.Errorf("failed to create db file: %w", err)
	}

	return nil
}

func RunMigrations(db *sql.DB, migrations embed.FS) error {
	log.Println("Running database migrations")
	d, err := iofs.New(migrations, "db/migrations")

	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		return fmt.Errorf("failed to instantiate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)

	if err != nil {
		return fmt.Errorf("failed to instantiate migrations: %w", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
