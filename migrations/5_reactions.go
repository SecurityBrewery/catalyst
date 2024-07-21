package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const ReactionCollectionName = "reactions"

func reactionsUp(db dbx.Builder) error {
	triggers := []string{"webhook", "hook"}
	reactions := []string{"python", "webhook"}

	return daos.New(db).SaveCollection(internalCollection(&models.Collection{
		Name: ReactionCollectionName,
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
			&schema.SchemaField{Name: "trigger", Type: schema.FieldTypeSelect, Required: true, Options: &schema.SelectOptions{MaxSelect: 1, Values: triggers}},
			&schema.SchemaField{Name: "triggerdata", Type: schema.FieldTypeJson, Required: true, Options: &schema.JsonOptions{MaxSize: 50_000}},
			&schema.SchemaField{Name: "action", Type: schema.FieldTypeSelect, Required: true, Options: &schema.SelectOptions{MaxSelect: 1, Values: reactions}},
			&schema.SchemaField{Name: "actiondata", Type: schema.FieldTypeJson, Required: true, Options: &schema.JsonOptions{MaxSize: 50_000}},
		),
	}))
}

func reactionsDown(db dbx.Builder) error {
	dao := daos.New(db)

	id, err := dao.FindCollectionByNameOrId(ReactionCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find collection %s: %w", ReactionCollectionName, err)
	}

	return dao.DeleteCollection(id)
}
