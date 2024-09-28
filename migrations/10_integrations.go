package migrations

import (
	"fmt"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const IntegrationCollectionName = "integrations"

func integrationsUp(db dbx.Builder) error {
	return daos.New(db).SaveCollection(&models.Collection{
		BaseModel: models.BaseModel{
			Id: IntegrationCollectionName,
		},
		Name: IntegrationCollectionName,
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
			&schema.SchemaField{Name: "plugin", Type: schema.FieldTypeText, Required: true},
			&schema.SchemaField{Name: "config", Type: schema.FieldTypeJson, Options: &schema.JsonOptions{MaxSize: 50_000}},
		),
	})
}

func integrationsDown(db dbx.Builder) error {
	dao := daos.New(db)

	id, err := dao.FindCollectionByNameOrId(IntegrationCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find collection %s: %w", IntegrationCollectionName, err)
	}

	return dao.DeleteCollection(id)
}
