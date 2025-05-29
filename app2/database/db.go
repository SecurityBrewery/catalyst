package database

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite" // import sqlite driver

	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
)

//go:embed migrations/*.sql
var schema embed.FS

func DB(filename string) (*sqlc.Queries, func(), error) {
	slog.Info("Connecting to database", "path", filename)

	_ = os.MkdirAll(filepath.Dir(filename), 0o755)

	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, nil, err
	}

	// Enabling Foreign Key Support
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, nil, err
	}

	return sqlc.New(db), func() {
		db.Close()
	}, nil
}

func migrateDB(db *sql.DB) error {
	iofsSource, err := iofs.New(schema, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs source: %w", err)
	}

	sqliteDriver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", iofsSource, "sqlite", sqliteDriver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
