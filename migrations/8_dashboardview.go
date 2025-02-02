package migrations

import (
	"github.com/pocketbase/pocketbase/core"
)

const dashboardCountsViewUpdateQuery = `SELECT id, count FROM (
  SELECT 'users' as id, COUNT(users.id) as count FROM users
  UNION
  SELECT 'tickets' as id, COUNT(tickets.id) as count FROM tickets
  UNION
  SELECT 'tasks' as id, COUNT(tasks.id) as count FROM tasks
  UNION 
  SELECT 'reactions' as id, COUNT(reactions.id) as count FROM reactions
) as counts;`

func dashboardCountsViewUpdateUp(app core.App) error {
	collection, err := app.FindCollectionByNameOrId(dashboardCountsViewName)
	if err != nil {
		return err
	}

	collection.ViewQuery = dashboardCountsViewUpdateQuery

	return app.Save(collection)
}

func dashboardCountsViewUpdateDown(app core.App) error {
	collection, err := app.FindCollectionByNameOrId(dashboardCountsViewName)
	if err != nil {
		return err
	}

	collection.ViewQuery = dashboardCountsViewQuery

	return app.Save(collection)
}
