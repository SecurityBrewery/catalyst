package migration

import (
	"context"
	"database/sql"
	"fmt"
)

func version(ctx context.Context, db *sql.DB) (int, error) {
	// get the current version of the database
	var currentVersion int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&currentVersion); err != nil {
		return 0, fmt.Errorf("failed to get current database version: %w", err)
	}

	return currentVersion, nil
}

func setVersion(ctx context.Context, db *sql.DB, version int) error {
	// Update the database version after successful migration
	_, err := db.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", version))
	if err != nil {
		return fmt.Errorf("failed to update database version: %w", err)
	}

	return nil
}
