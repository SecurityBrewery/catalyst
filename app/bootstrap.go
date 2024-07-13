package app

import (
	"log"

	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func Bootstrap(app core.App) error {
	if err := app.Bootstrap(); err != nil {
		return err
	}

	return migrateDBs(app)
}

func bootstrapCmd(app core.App) *cobra.Command {
	return &cobra.Command{
		Use: "bootstrap",
		Run: func(_ *cobra.Command, _ []string) {
			if err := Bootstrap(app); err != nil {
				log.Fatal(err)
			}
		},
	}
}
