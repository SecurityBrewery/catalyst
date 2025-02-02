package app

import (
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/migrations"
)

func HasFlag(app core.App, flag string) bool {
	records, err := app.FindAllRecords(migrations.FeatureCollectionName, dbx.HashExp{"name": flag})
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
	records, err := app.FindAllRecords(migrations.FeatureCollectionName)
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
	featureCollection, err := app.FindCollectionByNameOrId(migrations.FeatureCollectionName)
	if err != nil {
		return err
	}

	featureRecords, err := app.FindAllRecords(migrations.FeatureCollectionName)
	if err != nil {
		return err
	}

	var existingFlags []string //nolint:prealloc

	for _, featureRecord := range featureRecords {
		// remove feature flags that are not in the args
		if !slices.Contains(args, featureRecord.GetString("name")) {
			if err := app.Delete(featureRecord); err != nil {
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
		record := core.NewRecord(featureCollection)
		record.Set("name", arg)

		if err := app.Save(record); err != nil {
			return err
		}
	}

	return nil
}
