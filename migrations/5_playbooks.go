package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	PlaybookCollectionName = "playbooks"
	RunCollectionName      = "runs"
)

const dashboardCountsViewQuery2 = `SELECT id, count FROM (
  SELECT 'users' as id, COUNT(users.id) as count FROM users
  UNION
  SELECT 'tickets' as id, COUNT(tickets.id) as count FROM tickets
  UNION
  SELECT 'playbooks' as id, COUNT(playbooks.id) as count FROM playbooks
  UNION
  SELECT 'tasks' as id, COUNT(tasks.id) as count FROM tasks
) as counts;`

func playbookUp(db dbx.Builder) error {
	collections := []*models.Collection{
		internalCollection(&models.Collection{
			Name: PlaybookCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "steps", Type: schema.FieldTypeJson, Required: true, Options: &schema.JsonOptions{MaxSize: 50_000}},
			),
		}),
		internalCollection(&models.Collection{
			Name: RunCollectionName,
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{Name: "ticket", Type: schema.FieldTypeRelation, Required: true, Options: &schema.RelationOptions{CollectionId: TicketCollectionName, MaxSelect: types.Pointer(1), CascadeDelete: true}},
				&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true},
				&schema.SchemaField{Name: "steps", Type: schema.FieldTypeJson, Required: true, Options: &schema.JsonOptions{MaxSize: 50_000}},
			),
		}),
		internalView(dashboardCountsViewName, dashboardCountsViewQuery2),
	}

	dao := daos.New(db)
	for _, c := range collections {
		if err := dao.SaveCollection(c); err != nil {
			return err
		}
	}

	return nil
}

func playbookDown(db dbx.Builder) error {
	dao := daos.New(db)

	collections := []string{PlaybookCollectionName, RunCollectionName}

	for _, c := range collections {
		id, err := dao.FindCollectionByNameOrId(c)
		if err != nil {
			return err
		}

		if err := dao.DeleteCollection(id); err != nil {
			return err
		}
	}

	return nil
}
