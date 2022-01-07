package database

import (
	"context"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

type playbookResponse struct {
	PlaybookId   string         `json:"playbook_id"`
	PlaybookName string         `json:"playbook_name"`
	Playbook     model.Playbook `json:"playbook"`
	TicketId     int64          `json:"ticket_id"`
	TicketName   string         `json:"ticket_name"`
}

func (db *Database) TaskList(ctx context.Context) ([]*model.TaskWithContext, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `FOR d IN @@collection 
	` + ticketFilterQuery + `
	FILTER d.status == 'open'
	FOR playbook IN NOT_NULL(VALUES(d.playbooks), [])
	RETURN { ticket_id: TO_NUMBER(d._key), ticket_name: d.name, playbook_id: POSITION(d.playbooks, playbook, true), playbook_name: playbook.name, playbook: playbook }`
	cursor, _, err := db.Query(ctx, query, mergeMaps(ticketFilterVars, map[string]interface{}{
		"@collection": TicketCollectionName,
	}), busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var docs []*model.TaskWithContext
	for {
		var doc playbookResponse
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		playbook, err := toPlaybookResponse(&doc.Playbook)
		if err != nil {
			return nil, err
		}

		for _, task := range playbook.Tasks {
			if task.Active {
				docs = append(docs, &model.TaskWithContext{
					PlaybookId:   doc.PlaybookId,
					PlaybookName: doc.PlaybookName,
					Task:         task,
					TicketId:     doc.TicketId,
					TicketName:   doc.TicketName,
				})
			}
		}
	}

	return docs, err
}
