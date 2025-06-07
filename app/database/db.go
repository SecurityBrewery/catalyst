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
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// see https://briandouglas.ie/sqlite-defaults/ for more details

	pragmas := []string{
		// Enable WAL mode for better concurrency
		"journal_mode=WAL",
		// Enable synchronous mode for better data integrity
		"synchronous=NORMAL",
		// Set busy timeout to 5 seconds
		"busy_timeout=5000",
		// Set cache size to 20MB
		"cache_size=-20000",
		// Enable foreign key checks
		"foreign_keys=ON",
		// Enable incremental vacuuming
		"auto_vacuum=INCREMENTAL",
		// Set temp store to memory
		"temp_store=MEMORY",
		// Set mmap size to 2GB
		"mmap_size=2147483648",
		// Set page size to 8192
		"page_size=8192",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(fmt.Sprintf("PRAGMA %s", pragma)); err != nil {
			return nil, nil, fmt.Errorf("failed to set pragma %s: %w", pragma, err)
		}
	}

	if err := migrateDB(db); err != nil {
		return nil, nil, fmt.Errorf("failed to migrate database: %w", err)
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
