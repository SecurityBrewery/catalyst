package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

const ResourceCollectionName = "resources"

func resourcesUp(db dbx.Builder) error {
	return daos.New(db).SaveCollection(internalCollection(&models.Collection{
		Name: ResourceCollectionName,
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: TicketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
			&schema.SchemaField{Name: "service", Type: schema.FieldTypeText, Required: true},
			&schema.SchemaField{Name: "type", Type: schema.FieldTypeText, Required: true},
			&schema.SchemaField{Name: "resource", Type: schema.FieldTypeText, Required: true},
			&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
			&schema.SchemaField{Name: "icon", Type: schema.FieldTypeText, Required: true},
			&schema.SchemaField{Name: "description", Type: schema.FieldTypeText},
			&schema.SchemaField{Name: "url", Type: schema.FieldTypeText},
			&schema.SchemaField{Name: "attributes", Type: schema.FieldTypeJson, Options: &schema.JsonOptions{MaxSize: 50_000}},
		),
	}))
}

func resourcesDown(db dbx.Builder) error {
	dao := daos.New(db)

	id, err := dao.FindCollectionByNameOrId(ResourceCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find collection %s: %w", ResourceCollectionName, err)
	}

	return dao.DeleteCollection(id)
}
