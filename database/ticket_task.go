package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/arangodb/go-driver"
	"github.com/google/uuid"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/generated/time"
)

func (db *Database) TaskGet(ctx context.Context, id int64, playbookID string, taskID string) (*model.TicketWithTickets, *model.PlaybookResponse, *model.TaskWithContext, error) {
	inc, err := db.TicketGet(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}

	playbook, ok := inc.Playbooks[playbookID]
	if !ok {
		return nil, nil, nil, errors.New("playbook does not exist")
	}

	task, ok := playbook.Tasks[taskID]
	if !ok {
		return nil, nil, nil, errors.New("task does not exist")
	}

	return inc, playbook, &model.TaskWithContext{
		PlaybookId:   playbookID,
		PlaybookName: playbook.Name,
		TaskId:       taskID,
		Task:         task,
		TicketId:     id,
		TicketName:   inc.Name,
	}, nil
}

func (db *Database) TaskComplete(ctx context.Context, id int64, playbookID string, taskID string, data any) (*model.TicketWithTickets, error) {
	inc, err := db.TicketGet(ctx, id)
	if err != nil {
		return nil, err
	}

	completable := inc.Playbooks[playbookID].Tasks[taskID].Active
	if !completable {
		return nil, errors.New("cannot be completed")
	}

	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	LET playbook = d.playbooks[@playbookID]
	LET task = playbook.tasks[@taskID]
	LET newtask = MERGE(task, {"data": NOT_NULL(@data, {}), "done": true, closed: @closed })
	LET newtasks = MERGE(playbook.tasks, { @taskID: newtask } )
	LET newplaybook = MERGE(playbook, {"tasks": newtasks})
	LET newplaybooks = MERGE(d.playbooks, { @playbookID: newplaybook } )
	
	UPDATE d WITH { "modified": @now, "playbooks": newplaybooks } IN @@collection
	RETURN NEW`
	ticket, err := db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]any{
		"playbookID": playbookID,
		"taskID":     taskID,
		"data":       data,
		"closed":     time.Now().UTC(),
		"now":        time.Now().UTC(),
	}, ticketFilterVars), &busdb.Operation{
		Type: bus.DatabaseEntryUpdated,
		Ids: []driver.DocumentID{
			driver.NewDocumentID(TicketCollectionName, fmt.Sprintf("%d", id)),
		},
	})
	if err != nil {
		return nil, err
	}

	playbook := ticket.Playbooks[playbookID]
	task := playbook.Tasks[taskID]

	runNextTasks(id, playbookID, task.Next, task.Data, extractTicketResponse(ticket), db)

	return ticket, nil
}

func extractTicketResponse(ticket *model.TicketWithTickets) *model.TicketResponse {
	return &model.TicketResponse{
		Artifacts:  ticket.Artifacts,
		Comments:   ticket.Comments,
		Created:    ticket.Created,
		Details:    ticket.Details,
		Files:      ticket.Files,
		ID:         ticket.ID,
		Modified:   ticket.Modified,
		Name:       ticket.Name,
		Owner:      ticket.Owner,
		Playbooks:  ticket.Playbooks,
		Read:       ticket.Read,
		References: ticket.References,
		Schema:     ticket.Schema,
		Status:     ticket.Status,
		Type:       ticket.Type,
		Write:      ticket.Write,
	}
}

func (db *Database) TaskUpdateOwner(ctx context.Context, id int64, playbookID string, taskID string, owner string) (*model.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	LET playbook = d.playbooks[@playbookID]
	LET task = playbook.tasks[@taskID]
	LET newtask = MERGE(task, {"owner": @owner })
	LET newtasks = MERGE(playbook.tasks, { @taskID: newtask } )
	LET newplaybook = MERGE(playbook, {"tasks": newtasks})
	LET newplaybooks = MERGE(d.playbooks, { @playbookID: newplaybook } )
	
	UPDATE d WITH { "modified": @now, "playbooks": newplaybooks } IN @@collection
	RETURN NEW`
	ticket, err := db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]any{
		"playbookID": playbookID,
		"taskID":     taskID,
		"owner":      owner,
		"now":        time.Now().UTC(),
	}, ticketFilterVars), &busdb.Operation{
		Type: bus.DatabaseEntryUpdated,
		Ids: []driver.DocumentID{
			driver.NewDocumentID(TicketCollectionName, fmt.Sprintf("%d", id)),
		},
	})
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (db *Database) TaskUpdateData(ctx context.Context, id int64, playbookID string, taskID string, data map[string]any) (*model.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	LET playbook = d.playbooks[@playbookID]
	LET task = playbook.tasks[@taskID]
	LET newtask = MERGE(task, {"data": @data })
	LET newtasks = MERGE(playbook.tasks, { @taskID: newtask } )
	LET newplaybook = MERGE(playbook, {"tasks": newtasks})
	LET newplaybooks = MERGE(d.playbooks, { @playbookID: newplaybook } )
	
	UPDATE d WITH { "modified": @now, "playbooks": newplaybooks } IN @@collection
	RETURN NEW`
	ticket, err := db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]any{
		"playbookID": playbookID,
		"taskID":     taskID,
		"data":       data,
		"now":        time.Now().UTC(),
	}, ticketFilterVars), &busdb.Operation{
		Type: bus.DatabaseEntryUpdated,
		Ids: []driver.DocumentID{
			driver.NewDocumentID(TicketCollectionName, fmt.Sprintf("%d", id)),
		},
	})
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (db *Database) TaskRun(ctx context.Context, id int64, playbookID string, taskID string) error {
	ticket, _, task, err := db.TaskGet(ctx, id, playbookID, taskID)
	if err != nil {
		return err
	}

	if task.Task.Type == model.TaskTypeAutomation {
		if err := runTask(id, playbookID, taskID, task.Task, extractTicketResponse(ticket), db); err != nil {
			return err
		}
	}

	return nil
}

func runNextTasks(id int64, playbookID string, next map[string]string, data any, ticket *model.TicketResponse, db *Database) {
	for nextTaskID, requirement := range next {
		nextTask := ticket.Playbooks[playbookID].Tasks[nextTaskID]
		if nextTask.Type == model.TaskTypeAutomation {
			b, err := evalRequirement(requirement, data)
			if err != nil {
				continue
			}
			if b {
				if err := runTask(id, playbookID, nextTaskID, nextTask, ticket, db); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func runTask(ticketID int64, playbookID string, taskID string, task *model.TaskResponse, ticket *model.TicketResponse, db *Database) error {
	playbook := ticket.Playbooks[playbookID]
	msgContext := &model.Context{Playbook: playbook, Task: task, Ticket: ticket}
	origin := &model.Origin{TaskOrigin: &model.TaskOrigin{TaskId: taskID, PlaybookId: playbookID, TicketId: ticketID}}
	jobID := uuid.NewString()

	return publishJobMapping(jobID, *task.Automation, msgContext, origin, task.Payload, db)
}
