package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tools/types"
)

type baseUpFunc func(dao *daos.Dao) error

func baseUp(db dbx.Builder) error {
	dao := daos.New(db)

	for _, f := range []baseUpFunc{
		settingsUp,
		allowUserViewUp,
	} {
		if err := f(dao); err != nil {
			return err
		}
	}

	return nil
}

func settingsUp(dao *daos.Dao) error {
	s := settings.New()
	s.Meta.AppName = "Catalyst"
	s.Meta.HideControls = false

	return dao.SaveSettings(s)
}

func allowUserViewUp(dao *daos.Dao) error {
	collection, err := dao.FindCollectionByNameOrId(UserCollectionName)
	if err != nil {
		return err
	}

	collection.ViewRule = types.Pointer("@request.auth.id != ''")
	collection.ListRule = types.Pointer("@request.auth.id != ''")

	return dao.SaveCollection(collection)
}

func baseDown(db dbx.Builder) error {
	collection, err := daos.New(db).FindCollectionByNameOrId(UserCollectionName)
	if err != nil {
		return err
	}

	collection.ViewRule = types.Pointer("id = @request.auth.id")
	collection.ListRule = types.Pointer("id = @request.auth.id")

	return daos.New(db).SaveCollection(collection)
}
