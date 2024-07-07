package main

import (
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/migrations"
)

func flags(app core.App) ([]string, error) {
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
