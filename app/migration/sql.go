package migration

import (
	"context"
	"fmt"

	sqlmigrations "github.com/SecurityBrewery/catalyst/app/database/migrations"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/upload/uploader"
)

type sqlMigration struct {
	sqlName string
	upSQL   string
}

func newSQLMigration(name string) func() (migration, error) {
	return func() (migration, error) {
		up, err := sqlmigrations.Migrations.ReadFile(name + ".up.sql")
		if err != nil {
			return nil, fmt.Errorf("failed to read up migration file for %s: %w", name, err)
		}

		return &sqlMigration{
			sqlName: name,
			upSQL:   string(up),
		}, nil
	}
}

func (m sqlMigration) name() string {
	return m.sqlName
}

func (m sqlMigration) up(ctx context.Context, queries *sqlc.Queries, _ string, _ *uploader.Uploader) error {
	_, err := queries.WriteDB.ExecContext(ctx, m.upSQL)
	if err != nil {
		return fmt.Errorf("migration %s up failed: %w", m.sqlName, err)
	}

	return nil
}
