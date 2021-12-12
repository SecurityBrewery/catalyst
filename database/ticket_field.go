package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/iancoleman/strcase"
	"github.com/mingrammer/commonregex"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/pointer"
)

func (db *Database) AddArtifact(ctx context.Context, id int64, artifact *models.Artifact) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	if artifact.Status == nil {
		artifact.Status = pointer.String("unknown")
	}

	if artifact.Type == nil {
		artifact.Type = pointer.String(inferType(artifact.Name))
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	UPDATE d WITH { "modified": DATE_ISO8601(DATE_NOW()), "artifacts": PUSH(NOT_NULL(d.artifacts, []), @artifact) } IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{"artifact": artifact}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.DocumentID(fmt.Sprintf("%s/%d", TicketCollectionName, id)),
		},
		Msg: "Add artifact",
	})
}

func inferType(name string) string {
	switch {
	case commonregex.IPRegex.MatchString(name):
		return "ip"
	case commonregex.LinkRegex.MatchString(name):
		return "url"
	case commonregex.EmailRegex.MatchString(name):
		return "email"
	case commonregex.MD5HexRegex.MatchString(name):
		return "md5"
	case commonregex.SHA1HexRegex.MatchString(name):
		return "sha1"
	case commonregex.SHA256HexRegex.MatchString(name):
		return "sha256"
	}
	return "unknown"
}

func (db *Database) RemoveArtifact(ctx context.Context, id int64, name string) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	FOR a IN NOT_NULL(d.artifacts, [])
	FILTER a.name == @name
	LET newartifacts = REMOVE_VALUE(d.artifacts, a)
	UPDATE d WITH { "modified": DATE_ISO8601(DATE_NOW()), "artifacts": newartifacts } IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{"name": name}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.DocumentID(fmt.Sprintf("%s/%d", TicketCollectionName, id)),
		},
		Msg: "Remove artifact",
	})
}

func (db *Database) SetTemplate(ctx context.Context, id int64, schema string) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	UPDATE d WITH { "schema": @schema } IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{"schema": schema}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.DocumentID(fmt.Sprintf("%s/%d", TicketCollectionName, id)),
		},
		Msg: "Set Template",
	})
}

func (db *Database) AddComment(ctx context.Context, id int64, comment *models.CommentForm) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	if comment.Creator == nil || *comment.Creator == "" {
		user, exists := busdb.UserFromContext(ctx)
		if !exists {
			return nil, errors.New("no user in context")
		}

		comment.Creator = pointer.String(user.ID)
	}

	if comment.Created == nil {
		comment.Created = pointer.Time(time.Now().UTC())
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	UPDATE d WITH { "modified": DATE_ISO8601(DATE_NOW()), "comments": PUSH(NOT_NULL(d.comments, []), @comment) } IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{"comment": comment}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.DocumentID(fmt.Sprintf("%s/%d", TicketCollectionName, id)),
		},
		Msg: "Add comment",
	})
}

func (db *Database) RemoveComment(ctx context.Context, id int64, commentID int64) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	UPDATE d WITH { "modified": DATE_ISO8601(DATE_NOW()), "comments": REMOVE_NTH(d.comments, @commentID) } IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{"commentID": commentID}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.DocumentID(fmt.Sprintf("%s/%d", TicketCollectionName, id)),
		},
		Msg: "Remove comment",
	})
}

func (db *Database) SetReferences(ctx context.Context, id int64, references []*models.Reference) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	UPDATE d WITH { "modified": DATE_ISO8601(DATE_NOW()), "references": @references } IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{"references": references}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.DocumentID(fmt.Sprintf("%s/%d", TicketCollectionName, id)),
		},
		Msg: "Changed references",
	})
}

func (db *Database) LinkFiles(ctx context.Context, id int64, files []*models.File) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	UPDATE d WITH { "modified": DATE_ISO8601(DATE_NOW()), "files": @files } IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{"files": files}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.DocumentID(fmt.Sprintf("%s/%d", TicketCollectionName, id)),
		},
		Msg: "Linked files",
	})
}

func (db *Database) AddTicketPlaybook(ctx context.Context, id int64, playbookTemplate *models.PlaybookTemplateForm) (*models.TicketWithTickets, error) {
	pb, err := toPlaybook(playbookTemplate)
	if err != nil {
		return nil, err
	}

	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	playbookID := strcase.ToKebab(pb.Name)
	if playbookTemplate.ID != nil {
		playbookID = *playbookTemplate.ID
	}

	parentTicket, err := db.TicketGet(ctx, id)
	if err != nil {
		return nil, err
	}

	query := `FOR d IN @@collection 
	` + ticketFilterQuery + `
	FILTER d._key == @ID
	LET newplaybook =  ZIP( [@playbookID], [@playbook] )
	LET newplaybooks = MERGE(NOT_NULL(d.playbooks, {}), newplaybook)
	LET newticket = MERGE(d, { "modified": DATE_ISO8601(DATE_NOW()), "playbooks": newplaybooks })
	REPLACE d WITH newticket IN @@collection
	RETURN NEW`
	ticket, err := db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{
		"playbook":   pb,
		"playbookID": findName(parentTicket.Playbooks, playbookID),
	}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.NewDocumentID(TicketCollectionName, fmt.Sprintf("%d", id)),
		},
		Msg: "Added playbook",
	})
	if err != nil {
		return nil, err
	}

	if err := runRootTask(extractTicketResponse(ticket), playbookID, db); err != nil {
		return nil, err
	}

	return ticket, nil
}

func findName(playbooks map[string]*models.PlaybookResponse, name string) string {
	if _, ok := playbooks[name]; !ok {
		return name
	}

	for i := 0; ; i++ {
		try := fmt.Sprintf("%s%d", name, i)
		if _, ok := playbooks[try]; !ok {
			return try
		}
	}
}

func runRootTask(ticket *models.TicketResponse, playbookID string, db *Database) error {
	playbook := ticket.Playbooks[playbookID]

	var root *models.TaskResponse
	for _, task := range playbook.Tasks {
		if task.Order == 0 {
			root = task
		}
	}

	runNextTasks(ticket.ID, playbookID, root.Next, root.Data, ticket, db)
	return nil
}

func (db *Database) RemoveTicketPlaybook(ctx context.Context, id int64, playbookID string) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `FOR d IN @@collection 
	` + ticketFilterQuery + `
	FILTER d._key == @ID
	LET newplaybooks = UNSET(d.playbooks, @playbookID)
	REPLACE d WITH MERGE(d, { "modified": DATE_ISO8601(DATE_NOW()), "playbooks": newplaybooks }) IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{
		"playbookID": playbookID,
	}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.NewDocumentID(TicketCollectionName, fmt.Sprintf("%d", id)),
		},
		Msg: fmt.Sprintf("Removed playbook %s", playbookID),
	})
}
