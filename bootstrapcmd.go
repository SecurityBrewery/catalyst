package main

import (
	"log"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/migrations/logs"
	"github.com/pocketbase/pocketbase/tools/migrate"
	"github.com/spf13/cobra"
)

func bootstrapCmd(app *pocketbase.PocketBase) *cobra.Command {
	return &cobra.Command{
		Use: "bootstrap",
		Run: func(_ *cobra.Command, _ []string) {
			if err := app.Bootstrap(); err != nil {
				log.Fatal(err)
			}

			if err := migrateDBs(app); err != nil {
				log.Fatal(err)
			}
		},
	}
}

type migration struct {
	db         *dbx.DB
	migrations migrate.MigrationsList
}

func migrateDBs(app *pocketbase.PocketBase) error {
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
