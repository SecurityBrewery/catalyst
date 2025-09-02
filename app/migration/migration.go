package migration

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/upload"
)

type migration interface {
	name() string
	up(ctx context.Context, queries *sqlc.Queries, dir string, uploader *upload.Uploader) error
}

func Apply(ctx context.Context, queries *sqlc.Queries, dir string, uploader *upload.Uploader) error {
	currentVersion, err := version(ctx, queries.WriteDB)
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

		if err := m.up(ctx, queries, dir, uploader); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.name(), err)
		}
	}

	if err := setVersion(ctx, queries.WriteDB, currentVersion+len(migrations)); err != nil {
		return err
	}

	return nil
}
