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
	CommentCollectionName  = "comments"
	FeatureCollectionName  = "features"
	LinkCollectionName     = "links"
	TaskCollectionName     = "tasks"
	TicketCollectionName   = "tickets"
	TimelineCollectionName = "timeline"
	TypeCollectionName     = "types"
	WebhookCollectionName  = "webhooks"
	fileCollectionName     = "files"

	UserCollectionName = "_pb_users_auth_"
)

func collectionsUp(db dbx.Builder) error {
	collections := []*models.Collection{
		internalCollection(&models.Collection{
			Name: TypeCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "singular", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "plural", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "icon", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "schema", Type: schema.FieldTypeJson, Required: true, Options: &schema.JsonOptions{MaxSize: 50_000}},
			),
		}),
		internalCollection(&models.Collection{
			Name: TicketCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "type", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{CollectionId: TypeCollectionName, MaxSelect: types.Pointer(1)}},
				&schema.SchemaField{Name: "description", Type: schema.FieldTypeText},
				&schema.SchemaField{Name: "open", Type: schema.FieldTypeBool},
				&schema.SchemaField{Name: "resolution", Type: schema.FieldTypeText},
				&schema.SchemaField{Name: "schema", Type: schema.FieldTypeJson, Options: &schema.JsonOptions{MaxSize: 50_000}},
				&schema.SchemaField{Name: "state", Type: schema.FieldTypeJson, Options: &schema.JsonOptions{MaxSize: 50_000}},
				&schema.SchemaField{Name: "owner", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{CollectionId: UserCollectionName, MaxSelect: types.Pointer(1)}},
			),
		}),
		internalCollection(&models.Collection{
			Name: TaskCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: TicketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "open", Type: schema.FieldTypeBool},
				&schema.SchemaField{Name: "owner", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{CollectionId: UserCollectionName, MaxSelect: types.Pointer(1)}},
			),
		}),
		internalCollection(&models.Collection{
			Name: CommentCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: TicketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "author", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{CollectionId: UserCollectionName, MaxSelect: types.Pointer(1)}},
				&schema.SchemaField{Name: "message", Type: schema.FieldTypeText, Required: true},
			),
		}),
		internalCollection(&models.Collection{
			Name: TimelineCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: TicketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "time", Type: schema.FieldTypeDate, Required: true},
				&schema.SchemaField{Name: "message", Type: schema.FieldTypeText, Required: true},
			),
		}),
		internalCollection(&models.Collection{
			Name: LinkCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: TicketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "url", Type: schema.FieldTypeUrl, Required: true},
			),
		}),

		internalCollection(&models.Collection{
			Name: fileCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: TicketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "size", Type: schema.FieldTypeNumber, Required: true},
				&schema.SchemaField{Name: "blob", Type: schema.FieldTypeFile, Required: true, Options: &schema.FileOptions{MaxSelect: 1, MaxSize: 1024 * 1024 * 100}},
			),
		}),
		{
			BaseModel: models.BaseModel{
				Id: FeatureCollectionName,
			},
			Name: FeatureCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
			),
			ListRule: types.Pointer("@request.auth.id != ''"),
			ViewRule: types.Pointer("@request.auth.id != ''"),
			Indexes: types.JsonArray[string]{
				fmt.Sprintf("CREATE UNIQUE INDEX `unique_name` ON `%s` (`name`)", FeatureCollectionName),
			},
		},
	}

	dao := daos.New(db)
	for _, c := range collections {
		if err := dao.SaveCollection(c); err != nil {
			return err
		}
	}

	return nil
}

func internalCollection(c *models.Collection) *models.Collection {
	c.Id = c.Name
	c.ListRule = types.Pointer("@request.auth.id != ''")
	c.ViewRule = types.Pointer("@request.auth.id != ''")
	c.CreateRule = types.Pointer("@request.auth.id != ''")
	c.UpdateRule = types.Pointer("@request.auth.id != ''")
	c.DeleteRule = types.Pointer("@request.auth.id != ''")

	return c
}

func collectionsDown(db dbx.Builder) error {
	collections := []string{
		fileCollectionName,
		LinkCollectionName,
		TaskCollectionName,
		CommentCollectionName,
		TimelineCollectionName,
		FeatureCollectionName,
		TicketCollectionName,
		TypeCollectionName,
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
