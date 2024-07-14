package app

import (
	"slices"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/spf13/cobra"

	"github.com/SecurityBrewery/catalyst/migrations"
)

func Flags(app core.App) ([]string, error) {
	records, err := app.Dao().FindRecordsByExpr(migrations.FeatureCollectionName)
	if err != nil {
		return nil, err
	}

	var flags []string

	for _, r := range records {
		flags = append(flags, r.GetString("name"))
	}

	return flags, nil
}

func SetFlags(app core.App, args []string) error {
	featureCollection, err := app.Dao().FindCollectionByNameOrId(migrations.FeatureCollectionName)
	if err != nil {
		return err
	}

	featureRecords, err := app.Dao().FindRecordsByExpr(migrations.FeatureCollectionName)
	if err != nil {
		return err
	}

	var existingFlags []string

	for _, featureRecord := range featureRecords {
		// remove feature flags that are not in the args
		if !slices.Contains(args, featureRecord.GetString("name")) {
			if err := app.Dao().DeleteRecord(featureRecord); err != nil {
				return err
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
			return err
		}
	}

	return nil
}

func setFeatureFlagsCmd(app core.App) *cobra.Command {
	return &cobra.Command{
		Use: "set-feature-flags",
		Run: func(_ *cobra.Command, args []string) {
			if err := SetFlags(app, args); err != nil {
				app.Logger().Error(err.Error())
			}
		},
	}
}
