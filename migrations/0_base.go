package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type baseUpFunc func(app core.App) error

func baseUp(app core.App) error {
	for _, f := range []baseUpFunc{
		settingsUp,
		allowUserViewUp,
	} {
		if err := f(app); err != nil {
			return err
		}
	}

	return nil
}

func settingsUp(app core.App) error {
	s := app.Settings()
	s.Meta.AppName = "Catalyst"
	s.Meta.HideControls = false

	return app.Save(s)
}

func allowUserViewUp(app core.App) error {
	collection, err := app.FindCollectionByNameOrId(UserCollectionID)
	if err != nil {
		return err
	}

	collection.ViewRule = types.Pointer("@request.auth.id != ''")
	collection.ListRule = types.Pointer("@request.auth.id != ''")

	return app.Save(collection)
}

func baseDown(app core.App) error {
	collection, err := app.FindCollectionByNameOrId(UserCollectionID)
	if err != nil {
		return err
	}

	collection.ViewRule = types.Pointer("id = @request.auth.id")
	collection.ListRule = types.Pointer("id = @request.auth.id")

	return app.Save(collection)
}
