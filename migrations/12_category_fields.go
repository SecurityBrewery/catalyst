package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Adds a direct relation field `category` to tickets and tasks so users can pick a category
// directly from the Admin UI without creating a categorizations row manually.
func categoryFieldsUp(db dbx.Builder) error {
	dao := daos.New(db)

	// Helper to add field if missing
	addField := func(collectionName string) error {
		col, err := dao.FindCollectionByNameOrId(collectionName)
		if err != nil {
			return fmt.Errorf("failed to find collection %s: %w", collectionName, err)
		}

		if col.Schema.GetFieldByName("category") == nil {
			col.Schema.AddField(&schema.SchemaField{
				Name:     "category",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId: CategoryCollectionName,
					MaxSelect:    types.Pointer(1),
				},
			})
			return dao.SaveCollection(col)
		}
		return nil
	}

	if err := addField(TicketCollectionName); err != nil {
		return err
	}
	if err := addField(TaskCollectionName); err != nil {
		return err
	}

	return nil
}

func categoryFieldsDown(db dbx.Builder) error {
	dao := daos.New(db)

	removeField := func(collectionName string) error {
		col, err := dao.FindCollectionByNameOrId(collectionName)
		if err != nil {
			return fmt.Errorf("failed to find collection %s: %w", collectionName, err)
		}

		f := col.Schema.GetFieldByName("category")
		if f != nil {
			col.Schema.RemoveField(f.Id)
			return dao.SaveCollection(col)
		}
		return nil
	}

	if err := removeField(TicketCollectionName); err != nil {
		return err
	}
	if err := removeField(TaskCollectionName); err != nil {
		return err
	}
	return nil
}
