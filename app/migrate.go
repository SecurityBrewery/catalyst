package app

import (
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/migrations/logs"
	"github.com/pocketbase/pocketbase/tools/migrate"
)

type migration struct {
	db         *dbx.DB
	migrations migrate.MigrationsList
}

func MigrateDBs(app core.App) error {
	for _, m := range []migration{
		{db: app.DB(), migrations: migrations.AppMigrations},
		{db: app.LogsDB(), migrations: logs.LogsMigrations},
	} {
		runner, err := migrate.NewRunner(m.db, m.migrations)
		if err != nil {
			return err
		}

		if _, err := runner.Up(); err != nil {
			return err
		}
	}

	return nil
}

func isIgnored(err error) bool {
	// this fix ignores some errors that come from upstream migrations.
	ignoreErrors := []string{
		"1673167670_multi_match_migrate",
		"1660821103_add_user_ip_column",
	}

	for _, ignore := range ignoreErrors {
		if strings.Contains(err.Error(), ignore) {
			return true
		}
	}

	return false
}

func MigrateDBsDown(app core.App) error {
	for _, m := range []migration{
		{db: app.DB(), migrations: migrations.AppMigrations},
		{db: app.LogsDB(), migrations: logs.LogsMigrations},
	} {
		runner, err := migrate.NewRunner(m.db, m.migrations)
		if err != nil {
			return err
		}

		if _, err := runner.Down(len(m.migrations.Items())); err != nil {
			if isIgnored(err) {
				continue
			}

			return err
		}
	}

	return nil
}
