package migrations

import (
	"github.com/pocketbase/pocketbase/core"
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

func viewsUp(app core.App) error {
	collections := []*core.Collection{
		internalView(dashboardCountsViewName, dashboardCountsViewQuery),
		internalView(sidebarViewName, sidebarViewQuery),
	}

	for _, c := range collections {
		if err := app.Save(c); err != nil {
			return err
		}
	}

	return nil
}

func internalView(name, query string) *core.Collection {
	collection := core.NewViewCollection(name)
	collection.ViewQuery = query
	collection.ListRule = types.Pointer("@request.auth.id != ''")
	collection.ViewRule = types.Pointer("@request.auth.id != ''")

	return collection
}

func viewsDown(app core.App) error {
	collections := []string{dashboardCountsViewName, sidebarViewName}

	for _, c := range collections {
		id, err := app.FindCollectionByNameOrId(c)
		if err != nil {
			return err
		}

		if err := app.Delete(id); err != nil {
			return err
		}
	}

	return nil
}
