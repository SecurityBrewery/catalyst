package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

const searchViewName = "ticket_search"

const searchViewQuery = `
SELECT 
    tickets.id, 
    tickets.name,
    tickets.created,
    tickets.description,
    tickets.open,
    tickets.type,
    tickets.state,
    users.name as owner_name,
    group_concat(comments.message) as comment_messages,
    group_concat(files.name) as file_names,
    group_concat(links.name) as link_names,
    group_concat(links.url) as link_urls,
    group_concat(tasks.name) as task_names,
    group_concat(timeline.message) as timeline_messages
FROM tickets
LEFT JOIN comments ON comments.ticket = tickets.id
LEFT JOIN files ON files.ticket = tickets.id
LEFT JOIN links ON links.ticket = tickets.id
LEFT JOIN tasks ON tasks.ticket = tickets.id
LEFT JOIN timeline ON timeline.ticket = tickets.id
LEFT JOIN users ON users.id = tickets.owner
GROUP BY tickets.id
`

func searchViewUp(db dbx.Builder) error {
	return daos.New(db).SaveCollection(internalView(searchViewName, searchViewQuery))
}

func searchViewDown(db dbx.Builder) error {
	dao := daos.New(db)

	id, err := dao.FindCollectionByNameOrId(searchViewName)
	if err != nil {
		return err
	}

	return dao.DeleteCollection(id)
}
