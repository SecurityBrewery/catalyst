package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const ReactionCollectionName = "reactions"

func reactionsUp(db dbx.Builder) error {
	triggers := []string{"webhook", "hook"}
	reactions := []string{"python", "webhook"}

	collections := []*models.Collection{
		internalCollection(&models.Collection{
			Name: ReactionCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "trigger", Type: schema.FieldTypeSelect, Required: true, Options: &schema.SelectOptions{MaxSelect: 1, Values: triggers}},
				&schema.SchemaField{Name: "triggerdata", Type: schema.FieldTypeJson, Required: true},
				&schema.SchemaField{Name: "reaction", Type: schema.FieldTypeSelect, Required: true, Options: &schema.SelectOptions{MaxSelect: 1, Values: reactions}},
				&schema.SchemaField{Name: "reactiondata", Type: schema.FieldTypeJson, Required: true},
			),
		}),
	}

	dao := daos.New(db)
	for _, c := range collections {
		if err := dao.SaveCollection(c); err != nil {
			return err
		}
	}

	return nil
}

func reactionsDown(db dbx.Builder) error {
	collections := []string{
		ReactionCollectionName,
	}

	dao := daos.New(db)

	for _, name := range collections {
		id, err := dao.FindCollectionByNameOrId(name)
		if err != nil {
			return err
		}

		if err := dao.DeleteCollection(id); err != nil {
			return err
		}
	}

	return nil
}
