package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	dashboardCountsViewName = "dashboard_counts"
	sidebarViewName         = "sidebar"
)

const dashboardCountsViewQuery = `SELECT id, count FROM (
  SELECT 'users' as id, COUNT(users.id) as count FROM users
  UNION
  SELECT 'tickets' as id, COUNT(tickets.id) as count FROM tickets
  UNION
  SELECT 'tasks' as id, COUNT(tasks.id) as count FROM tasks
) as counts;`

const sidebarViewQuery = `SELECT types.id as id, types.singular as singular, types.plural as plural, types.icon as icon, (SELECT COUNT(tickets.id) FROM tickets WHERE tickets.type = types.id AND tickets.open = true) as count
FROM types
ORDER BY types.plural;`

func viewsUp(db dbx.Builder) error {
	collections := []*models.Collection{
		internalView(dashboardCountsViewName, dashboardCountsViewQuery),
		internalView(sidebarViewName, sidebarViewQuery),
	}

	dao := daos.New(db)
	for _, c := range collections {
		if err := dao.SaveCollection(c); err != nil {
			return err
		}
	}

	return nil
}

func internalView(name, query string) *models.Collection {
	return &models.Collection{
		Name:     name,
		Type:     models.CollectionTypeView,
		Options:  types.JsonMap{"query": query},
		ListRule: types.Pointer("@request.auth.id != ''"),
		ViewRule: types.Pointer("@request.auth.id != ''"),
	}
}

func viewsDown(db dbx.Builder) error {
	dao := daos.New(db)

	collections := []string{dashboardCountsViewName, sidebarViewName}

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
