package main

import (
	"log"
	"slices"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/spf13/cobra"

	"github.com/SecurityBrewery/catalyst/fakedata"
	"github.com/SecurityBrewery/catalyst/migrations"
)

func fakeDataCmd(app *pocketbase.PocketBase) *cobra.Command {
	var userCount, ticketCount int

	cmd := &cobra.Command{
		Use: "fake-data",
		Run: func(_ *cobra.Command, _ []string) {
			if err := fakedata.Generate(app.DB(), userCount, ticketCount); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.PersistentFlags().IntVar(&userCount, "users", 10, "Number of users to generate")

	cmd.PersistentFlags().IntVar(&ticketCount, "tickets", 100, "Number of tickets to generate")

	return cmd
}

func setFeatureFlagsCmd(app *pocketbase.PocketBase) *cobra.Command {
	return &cobra.Command{
		Use: "set-feature-flags",
		Run: func(_ *cobra.Command, args []string) {
			featureCollection, err := app.Dao().FindCollectionByNameOrId(migrations.FeatureCollectionName)
			if err != nil {
				log.Fatal(err)
			}

			featureRecords, err := app.Dao().FindRecordsByExpr(migrations.FeatureCollectionName)
			if err != nil {
				log.Fatal(err)
			}

			var existingFlags []string

			for _, featureRecord := range featureRecords {
				// remove feature flags that are not in the args
				if !slices.Contains(args, featureRecord.GetString("name")) {
					if err := app.Dao().DeleteRecord(featureRecord); err != nil {
						log.Fatal(err)
					}

					continue
				}

				existingFlags = append(existingFlags, featureRecord.GetString("name"))
			}

			for _, arg := range args {
				if slices.Contains(existingFlags, arg) {
					continue
				}

				// add feature flags that are not in the args
				record := models.NewRecord(featureCollection)
				record.Set("name", arg)

				if err := app.Dao().SaveRecord(record); err != nil {
					log.Fatal(err)
				}
			}
		},
	}
}
