package database

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

func (db *Database) Settings(ctx context.Context) (*model.Settings, error) {
	settings := &model.Settings{}
	if _, err := db.settingsCollection.ReadDocument(ctx, "global", settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func (db *Database) SaveSettings(ctx context.Context, settings *model.Settings) (*model.Settings, error) {
	exists, err := db.settingsCollection.DocumentExists(ctx, "global")
	if err != nil {
		return nil, err
	}

	if exists {
		if _, err := db.settingsCollection.ReplaceDocument(ctx, "global", settings); err != nil {
			return nil, err
		}
	} else {
		if _, err := db.settingsCollection.CreateDocument(ctx, ctx, "global", settings); err != nil {
			return nil, err
		}
	}

	return settings, nil
}
