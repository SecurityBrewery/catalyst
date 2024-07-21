package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const SystemUserID = "system"

func systemuserUp(db dbx.Builder) error {
	dao := daos.New(db)

	collection, err := dao.FindCollectionByNameOrId(UserCollectionName)
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)
	record.SetId(SystemUserID)
	record.Set("name", "system")
	record.Set("username", "system")
	record.Set("verified", true)

	return dao.SaveRecord(record)
}

func systemuserDown(db dbx.Builder) error {
	dao := daos.New(db)

	record, err := dao.FindRecordById(UserCollectionName, SystemUserID)
	if err != nil {
		return err
	}

	return dao.DeleteRecord(record)
}
