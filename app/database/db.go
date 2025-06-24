package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3" // import sqlite driver
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/migration"
)

const sqliteDriver = "sqlite3"

func DB(ctx context.Context, filename string) (*sqlc.Queries, func(), error) {
	slog.InfoContext(ctx, "Connecting to database", "path", filename)

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

	_ = os.MkdirAll(filepath.Dir(filename), 0o755)

	write, err := sql.Open(sqliteDriver, fmt.Sprintf("file:%s", filename))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	write.SetMaxOpenConns(1)
	write.SetConnMaxIdleTime(time.Minute)

	for _, pragma := range pragmas {
		if _, err := write.ExecContext(ctx, fmt.Sprintf("PRAGMA %s", pragma)); err != nil {
			return nil, nil, fmt.Errorf("failed to set pragma %s: %w", pragma, err)
		}
	}

	read, err := sql.Open(sqliteDriver, fmt.Sprintf("file:%s?mode=ro", filename))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	read.SetMaxOpenConns(100)
	read.SetConnMaxIdleTime(time.Minute)

	if err := migration.Apply(ctx, write); err != nil {
		return nil, nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return sqlc.New(read, write), func() {
		if err := read.Close(); err != nil {
			slog.Error("failed to close read connection", "error", err)
		}

		if err := write.Close(); err != nil {
			slog.Error("failed to close write connection", "error", err)
		}
	}, nil
}

func TestDB(t *testing.T, dir string) *sqlc.Queries {
	queries, cleanup, err := DB(t.Context(), filepath.Join(dir, "data.db"))
	require.NoError(t, err)
	t.Cleanup(cleanup)

	return queries
}

func GenerateID(prefix string) string {
	return strings.ToLower(prefix) + randomstring(12)
}

const base32alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func randomstring(l int) string {
	rand.Text()

	src := make([]byte, l)
	_, _ = rand.Read(src)

	for i := range src {
		src[i] = base32alphabet[int(src[i])%len(base32alphabet)]
	}

	return string(src)
}
