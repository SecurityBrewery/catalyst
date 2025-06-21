package migration

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type migration interface {
	name() string
	up(ctx context.Context, db *sql.DB) error
	down(ctx context.Context, db *sql.DB) error
}

func Apply(ctx context.Context, db *sql.DB) error {
	currentVersion, err := version(ctx, db)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Current database version", "version", currentVersion)

	migrations, err := migrations(currentVersion)
	if err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}

	if len(migrations) == 0 {
		slog.InfoContext(ctx, "No migrations to apply")

		return nil
	}

	for _, m := range migrations {
		slog.InfoContext(ctx, "Applying migration", "name", m.name())

		if err := m.up(ctx, db); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.name(), err)
		}
	}

	if err := setVersion(ctx, db, currentVersion+len(migrations)); err != nil {
		return err
	}

	return nil
}

func Rollback(ctx context.Context, db *sql.DB) error {
	currentVersion, err := version(ctx, db)
	if err != nil {
		return err
	}

	if currentVersion == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	migrations, err := migrations(0)
	if err != nil {
		return fmt.Errorf("failed to get migrations for rollback: %w", err)
	}

	if currentVersion > len(migrations) {
		return fmt.Errorf("current version %d exceeds available migrations %d", currentVersion, len(migrations))
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		slog.InfoContext(ctx, "Rolling back migration", "name", m.name())

		if err := m.down(ctx, db); err != nil {
			return fmt.Errorf("rollback %s failed: %w", m.name(), err)
		}
	}

	if err := setVersion(ctx, db, 0); err != nil {
		return fmt.Errorf("failed to update version after rollback: %w", err)
	}

	return nil
}
