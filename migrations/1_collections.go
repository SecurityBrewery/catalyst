package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"

	"github.com/SecurityBrewery/catalyst/ff"
)

const (
	timelineCollectionName = "timeline"
	commentCollectionName  = "comments"
	fileCollectionName     = "files"
	linkCollectionName     = "links"
	playbookCollectionName = "playbooks"
	runCollectionName      = "runs"
	taskCollectionName     = "tasks"
	ticketCollectionName   = "tickets"
	typeCollectionName     = "types"
	webhookCollectionName  = "webhooks"

	userCollectionName = "_pb_users_auth_"
)

func collectionsUp(db dbx.Builder) error {
	collections := []*models.Collection{
		internalCollection(&models.Collection{
			Name: typeCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "singular", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "plural", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "icon", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "schema", Type: schema.FieldTypeJson, Required: true, Options: &schema.JsonOptions{MaxSize: 50_000}},
			),
		}),
		internalCollection(&models.Collection{
			Name: ticketCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "type", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{CollectionId: typeCollectionName, MaxSelect: types.Pointer(1)}},
				&schema.SchemaField{Name: "description", Type: schema.FieldTypeText},
				&schema.SchemaField{Name: "open", Type: schema.FieldTypeBool},
				&schema.SchemaField{Name: "schema", Type: schema.FieldTypeJson, Options: &schema.JsonOptions{MaxSize: 50_000}},
				&schema.SchemaField{Name: "state", Type: schema.FieldTypeJson, Options: &schema.JsonOptions{MaxSize: 50_000}},
				&schema.SchemaField{Name: "owner", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{CollectionId: userCollectionName, MaxSelect: types.Pointer(1)}},
			),
		}),
		internalCollection(&models.Collection{
			Name: taskCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: ticketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "open", Type: schema.FieldTypeBool},
				&schema.SchemaField{Name: "owner", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{CollectionId: userCollectionName, MaxSelect: types.Pointer(1)}},
			),
		}),
		internalCollection(&models.Collection{
			Name: commentCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: ticketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "author", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{CollectionId: userCollectionName, MaxSelect: types.Pointer(1)}},
				&schema.SchemaField{Name: "message", Type: schema.FieldTypeText, Required: true},
			),
		}),
		internalCollection(&models.Collection{
			Name: timelineCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: ticketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "time", Type: schema.FieldTypeDate, Required: true},
				&schema.SchemaField{Name: "message", Type: schema.FieldTypeText, Required: true},
			),
		}),
		internalCollection(&models.Collection{
			Name: linkCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: ticketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "url", Type: schema.FieldTypeUrl, Required: true},
			),
		}),

		internalCollection(&models.Collection{
			Name: playbookCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "steps", Type: schema.FieldTypeJson, Required: true, Options: &schema.JsonOptions{MaxSize: 50_000}},
			),
		}),
		internalCollection(&models.Collection{
			Name: runCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: ticketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "steps", Type: schema.FieldTypeJson, Required: true, Options: &schema.JsonOptions{MaxSize: 50_000}},
			),
		}),
	}

	if !ff.HasDemoFlag() {
		collections = append(collections,
			internalCollection(&models.Collection{
				Name: fileCollectionName,
				Type: models.CollectionTypeBase,
				Schema: schema.NewSchema(
					&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: ticketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
					&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
					&schema.SchemaField{Name: "size", Type: schema.FieldTypeNumber, Required: true},
					&schema.SchemaField{Name: "blob", Type: schema.FieldTypeFile, Required: true, Options: &schema.FileOptions{MaxSelect: 1, MaxSize: 1024 * 1024 * 100}},
				),
			}),
		)
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
		playbookCollectionName,
		ticketCollectionName,
		typeCollectionName,
		fileCollectionName,
		linkCollectionName,
		taskCollectionName,
		runCollectionName,
		commentCollectionName,
		timelineCollectionName,
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
