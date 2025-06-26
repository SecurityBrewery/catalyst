package migration

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sqlmigrations "github.com/SecurityBrewery/catalyst/app/database/migrations"
)

type sqlMigration struct {
	sqlName string
	upSQL   string
	downSQL string
}

func newSQLMigration(name string) func() (migration, error) {
	return func() (migration, error) {
		up, err := sqlmigrations.Migrations.ReadFile(name + ".up.sql")
		if err != nil {
			return nil, fmt.Errorf("failed to read up migration file for %s: %w", name, err)
		}

		down, err := sqlmigrations.Migrations.ReadFile(name + ".down.sql")
		if err != nil {
			return nil, fmt.Errorf("failed to read down migration file for %s: %w", name, err)
		}

		return &sqlMigration{
			sqlName: name,
			upSQL:   string(up),
			downSQL: string(down),
		}, nil
	}
}

func (m sqlMigration) name() string {
	return m.sqlName
}

func (m sqlMigration) up(ctx context.Context, db *sql.DB) error {
	parts := strings.Split(m.upSQL, ";")

	for _, part := range parts {
		_, err := db.ExecContext(ctx, part)
		if err != nil {
			return fmt.Errorf("migration %s up failed: %w", part, err)
		}
	}

	return nil
}

func (m sqlMigration) down(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, m.downSQL)
	if err != nil {
		return fmt.Errorf("migration %s down failed: %w", m.sqlName, err)
	}

	return nil
}
