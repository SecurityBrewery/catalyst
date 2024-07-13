package app

import (
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

func migrateDBs(app core.App) error {
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
