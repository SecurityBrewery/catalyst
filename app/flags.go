package app

import (
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/spf13/cobra"

	"github.com/SecurityBrewery/catalyst/migrations"
)

func HasFlag(app core.App, flag string) bool {
	records, err := app.Dao().FindRecordsByExpr(migrations.FeatureCollectionName, dbx.HashExp{"name": flag})
	if err != nil {
		app.Logger().Error(err.Error())

		return false
	}

	for _, r := range records {
		if r.GetString("name") == flag {
			return true
		}
	}

	return false
}

func Flags(app core.App) ([]string, error) {
	records, err := app.Dao().FindRecordsByExpr(migrations.FeatureCollectionName)
	if err != nil {
		return nil, err
	}

	flags := make([]string, 0, len(records))

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

	var existingFlags []string //nolint:prealloc

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

func setAppURL(app core.App) *cobra.Command {
	return &cobra.Command{
		Use: "set-app-url",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) != 1 {
				app.Logger().Error("missing app url")

				return
			}

			settings := app.Settings()

			settings.Meta.AppUrl = args[0]

			if err := app.Dao().SaveSettings(settings); err != nil {
				app.Logger().Error(err.Error())
			}
		},
	}
}
