package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/types"
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

func dashboardCountsViewUpdateUp(db dbx.Builder) error {
	dao := daos.New(db)

	collection, err := dao.FindCollectionByNameOrId(dashboardCountsViewName)
	if err != nil {
		return err
	}

	collection.Options = types.JsonMap{"query": dashboardCountsViewUpdateQuery}

	return dao.SaveCollection(collection)
}

func dashboardCountsViewUpdateDown(db dbx.Builder) error {
	dao := daos.New(db)

	collection, err := dao.FindCollectionByNameOrId(dashboardCountsViewName)
	if err != nil {
		return err
	}

	collection.Options = types.JsonMap{"query": dashboardCountsViewQuery}

	return dao.SaveCollection(collection)
}
