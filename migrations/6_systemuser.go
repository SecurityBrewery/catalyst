package migrations

import (
	"github.com/pocketbase/pocketbase/core"
)

const SystemUserID = "system"

func systemuserUp(app core.App) error {
	collection, err := app.FindCollectionByNameOrId(UserCollectionID)
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	record.Id = SystemUserID
	record.Set("name", "system")
	record.Set("username", "system")
	record.Set("verified", true)

	return app.SaveNoValidate(record)
}

func systemuserDown(app core.App) error {
	record, err := app.FindRecordById(UserCollectionID, SystemUserID)
	if err != nil {
		return err
	}

	return app.Delete(record)
}
