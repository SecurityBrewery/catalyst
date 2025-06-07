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

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
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

	// see https://briandouglas.ie/sqlite-defaults/ for more details

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, nil, err
	}

	// Enable synchronous mode for better data integrity
	if _, err := db.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		return nil, nil, err
	}

	// Set busy timeout to 5 seconds
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		return nil, nil, err
	}

	// Set cache size to 20MB
	if _, err := db.Exec("PRAGMA cache_size = -20000"); err != nil {
		return nil, nil, err
	}

	// Enable foreign key checks
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, nil, err
	}

	// Enable incremental vacuuming
	if _, err := db.Exec("PRAGMA auto_vacuum = INCREMENTAL"); err != nil {
		return nil, nil, err
	}

	// Set temp store to memory
	if _, err := db.Exec("PRAGMA temp_store = MEMORY"); err != nil {
		return nil, nil, err
	}

	// Set mmap size to 2GB
	if _, err := db.Exec("PRAGMA mmap_size = 2147483648"); err != nil {
		return nil, nil, err
	}

	// Set page size to 8192
	if _, err := db.Exec("PRAGMA page_size = 8192"); err != nil {
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
