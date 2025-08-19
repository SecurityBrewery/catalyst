package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const (
	CategoryCollectionName = "categories"
)

func categoriesUp(db dbx.Builder) error {
	collection := &models.Collection{
		Name: CategoryCollectionName,
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "title", Type: schema.FieldTypeText, Required: true},
		),
	}

	collection = internalCollection(collection)

	return daos.New(db).SaveCollection(collection)
}

func categoriesDown(db dbx.Builder) error {
	dao := daos.New(db)

	id, err := dao.FindCollectionByNameOrId(CategoryCollectionName)
	if err != nil {
		return err
	}

	return dao.DeleteCollection(id)
}