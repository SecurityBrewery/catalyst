package app

import (
	"github.com/pocketbase/pocketbase/core"
)

func MigrateDBs(app core.App) error {
	return app.RunAllMigrations()
}

func MigrateDBsDown(app core.App) error {
	for _, m := range []core.MigrationsList{
		core.AppMigrations,
		core.SystemMigrations,
	} {
		if _, err := core.NewMigrationsRunner(app, m).Down(len(m.Items())); err != nil {
			return err
		}
	}

	return nil
}
