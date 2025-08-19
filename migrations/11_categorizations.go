package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	CategorizationCollectionName = "categorizations"
)

func categorizationsUp(db dbx.Builder) error {
	collection := &models.Collection{
		Name: CategorizationCollectionName,
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "category", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: CategoryCollectionName, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Name: "target_collection", Type: schema.FieldTypeText, Required: true},
			&schema.SchemaField{Name: "target_id", Type: schema.FieldTypeText, Required: true},
		),
		Indexes: types.JsonArray[string]{
			fmt.Sprintf("CREATE UNIQUE INDEX `unique_target_category` ON `%s` (`target_collection`, `target_id`, `category`)", CategorizationCollectionName),
		},
	}

	collection = internalCollection(collection)

	return daos.New(db).SaveCollection(collection)
}

func categorizationsDown(db dbx.Builder) error {
	dao := daos.New(db)

	id, err := dao.FindCollectionByNameOrId(CategorizationCollectionName)
	if err != nil {
		return err
	}

	return dao.DeleteCollection(id)
}
