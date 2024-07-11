package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const (
	TriggerWebhookCollectionName  = "triggers_webhooks"
	TriggerHookCollectionName     = "triggers_hooks"
	ReactionViewName              = "reactions"
	ReactionPythonCollectionName  = "reactions_python"
	ReactionWebhookCollectionName = "reactions_webhooks"
)

const reactionViewQuery = `SELECT id, name, type, created, updated FROM (
  SELECT id, name, created, updated, 'python' as type FROM reactions_python
  UNION
  SELECT id, name, created, updated, 'webhook' as type FROM reactions_webhooks
) as reactions;`

func actionsUp(db dbx.Builder) error {
	hookCollections := []string{TicketCollectionName, TaskCollectionName, CommentCollectionName, TimelineCollectionName, LinkCollectionName, fileCollectionName}
	hookEvents := []string{"create", "update", "delete"}

	collections := []*models.Collection{
		internalCollection(&models.Collection{
			Name: TriggerWebhookCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "token", Type: schema.FieldTypeText},
				&schema.SchemaField{Name: "path", Type: schema.FieldTypeText, Required: true},
			),
		}),
		internalCollection(&models.Collection{
			Name: TriggerHookCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "collection", Type: schema.FieldTypeSelect, Required: true, Options: &schema.SelectOptions{Values: hookCollections}},
				&schema.SchemaField{Name: "event", Type: schema.FieldTypeSelect, Required: true, Options: &schema.SelectOptions{Values: hookEvents}},
			),
		}),
		internalCollection(&models.Collection{
			Name: ReactionPythonCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "requirements", Type: schema.FieldTypeText},
				&schema.SchemaField{Name: "script", Type: schema.FieldTypeText, Required: true},
			),
		}),
		internalCollection(&models.Collection{
			Name: ReactionWebhookCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "headers", Type: schema.FieldTypeText},
				&schema.SchemaField{Name: "destination", Type: schema.FieldTypeText, Required: true},
			),
		}),
		internalView(ReactionViewName, reactionViewQuery),
	}

	dao := daos.New(db)
	for _, c := range collections {
		if err := dao.SaveCollection(c); err != nil {
			return err
		}
	}

	return nil
}

func actionsDown(db dbx.Builder) error {
	collections := []string{
		TriggerWebhookCollectionName,
		TriggerHookCollectionName,
		ReactionPythonCollectionName,
		ReactionWebhookCollectionName,
		ReactionViewName,
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
