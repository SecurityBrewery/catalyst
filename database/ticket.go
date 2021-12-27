package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/arangodb/go-driver"
	"github.com/xeipuuv/gojsonschema"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/caql"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/index"
	"github.com/SecurityBrewery/catalyst/time"
)

func toTicket(ticketForm *models.TicketForm) (interface{}, error) {
	playbooks, err := toPlaybooks(ticketForm.Playbooks)
	if err != nil {
		return nil, err
	}

	ticket := &models.Ticket{
		Artifacts:  ticketForm.Artifacts,
		Comments:   ticketForm.Comments,
		Details:    ticketForm.Details,
		Files:      ticketForm.Files,
		Name:       ticketForm.Name,
		Owner:      ticketForm.Owner,
		Playbooks:  playbooks,
		Read:       ticketForm.Read,
		References: ticketForm.References,
		Status:     ticketForm.Status,
		Type:       ticketForm.Type,
		Write:      ticketForm.Write,
		// ID:         ticketForm.ID,
		// Created:    ticketForm.Created,
		// Modified:   ticketForm.Modified,
		// Schema:     ticketForm.Schema,
	}

	if ticketForm.Created != nil {
		ticket.Created = *ticketForm.Created
	} else {
		ticket.Created = time.Now().UTC()
	}
	if ticketForm.Modified != nil {
		ticket.Modified = *ticketForm.Modified
	} else {
		ticket.Modified = time.Now().UTC()
	}
	if ticketForm.Schema != nil {
		ticket.Schema = *ticketForm.Schema
	} else {
		ticket.Schema = "{}"
	}
	if ticketForm.Status == "" {
		ticket.Status = "open"
	}
	if ticketForm.ID != nil {
		return &busdb.Keyed{Key: strconv.FormatInt(*ticketForm.ID, 10), Doc: ticket}, nil
	}
	return ticket, nil
}

func toTicketResponses(tickets []*models.TicketSimpleResponse) ([]*models.TicketResponse, error) {
	var extendedTickets []*models.TicketResponse
	for _, simple := range tickets {
		tr, err := toTicketResponse(simple)
		if err != nil {
			return nil, err
		}
		extendedTickets = append(extendedTickets, tr)
	}
	return extendedTickets, nil
}

func toTicketResponse(ticket *models.TicketSimpleResponse) (*models.TicketResponse, error) {
	playbooks, err := toPlaybookResponses(ticket.Playbooks)
	if err != nil {
		return nil, err
	}

	return &models.TicketResponse{
		ID:         ticket.ID,
		Artifacts:  ticket.Artifacts,
		Comments:   ticket.Comments,
		Created:    ticket.Created,
		Details:    ticket.Details,
		Files:      ticket.Files,
		Modified:   ticket.Modified,
		Name:       ticket.Name,
		Owner:      ticket.Owner,
		Playbooks:  playbooks,
		Read:       ticket.Read,
		References: ticket.References,
		Schema:     ticket.Schema,
		Status:     ticket.Status,
		Type:       ticket.Type,
		Write:      ticket.Write,
	}, nil
}

func toTicketSimpleResponse(key string, ticket *models.Ticket) (*models.TicketSimpleResponse, error) {
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return nil, err
	}

	return &models.TicketSimpleResponse{
		Artifacts:  ticket.Artifacts,
		Comments:   ticket.Comments,
		Created:    ticket.Created,
		Details:    ticket.Details,
		Files:      ticket.Files,
		ID:         id,
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
	}, nil
}

func toTicketWithTickets(ticketResponse *models.TicketResponse, tickets []*models.TicketSimpleResponse) *models.TicketWithTickets {
	return &models.TicketWithTickets{
		Artifacts:  ticketResponse.Artifacts,
		Comments:   ticketResponse.Comments,
		Created:    ticketResponse.Created,
		Details:    ticketResponse.Details,
		Files:      ticketResponse.Files,
		ID:         ticketResponse.ID,
		Modified:   ticketResponse.Modified,
		Name:       ticketResponse.Name,
		Owner:      ticketResponse.Owner,
		Playbooks:  ticketResponse.Playbooks,
		Read:       ticketResponse.Read,
		References: ticketResponse.References,
		Schema:     ticketResponse.Schema,
		Status:     ticketResponse.Status,
		Type:       ticketResponse.Type,
		Write:      ticketResponse.Write,

		Tickets: tickets,
	}
}

func toPlaybookResponses(playbooks map[string]*models.Playbook) (map[string]*models.PlaybookResponse, error) {
	pr := map[string]*models.PlaybookResponse{}
	var err error
	for k, v := range playbooks {
		pr[k], err = toPlaybookResponse(v)
		if err != nil {
			return nil, err
		}
	}
	return pr, nil
}

func toPlaybookResponse(playbook *models.Playbook) (*models.PlaybookResponse, error) {
	graph, err := playbookGraph(playbook)
	if err != nil {
		return nil, err
	}

	re := &models.PlaybookResponse{
		Name:  playbook.Name,
		Tasks: map[string]*models.TaskResponse{},
	}

	results, err := graph.Toposort()
	if err != nil {
		return nil, err
	}

	i := 0
	for _, taskID := range results {
		rootTask, err := toTaskResponse(playbook, taskID, i, graph)
		if err != nil {
			return nil, err
		}
		re.Tasks[taskID] = rootTask
		i++
	}
	return re, nil
}

func (db *Database) TicketBatchCreate(ctx context.Context, ticketForms []*models.TicketForm) ([]*models.TicketResponse, error) {
	update, err := db.Hooks.IngestionFilter(ctx, db.Index)
	if err != nil {
		return nil, err
	}

	var dbTickets []interface{}
	for _, ticketForm := range ticketForms {
		ticket, err := toTicket(ticketForm)
		if err != nil {
			return nil, err
		}

		if err := validate(ticket, models.TicketSchema); err != nil {
			return nil, err
		}

		dbTickets = append(dbTickets, ticket)
	}

	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `FOR d IN @tickets
	` + ticketFilterQuery + `
	LET updates = ` + update + `
	LET newdoc = LENGTH(updates) != 0 ? APPLY("MERGE_RECURSIVE", APPEND([d], updates)) : d
	LET keyeddoc = HAS(newdoc, "id") ? MERGE(newdoc, {"_key": TO_STRING(newdoc.id)}) : newdoc
	LET noiddoc = UNSET(keyeddoc, "id")
	INSERT noiddoc INTO @@collection
	RETURN NEW`
	apiTickets, _, err := db.ticketListQuery(ctx, query, mergeMaps(map[string]interface{}{
		"tickets": dbTickets,
	}, ticketFilterVars), busdb.CreateOperation)
	if err != nil {
		return nil, err
	}

	if err = batchIndex(db.Index, apiTickets); err != nil {
		return nil, err
	}

	var ids []driver.DocumentID
	for _, apiTicket := range apiTickets {
		ids = append(ids, driver.NewDocumentID(TicketCollectionName, fmt.Sprint(apiTicket.ID)))
	}

	go db.bus.PublishDatabaseUpdate(ids, bus.DatabaseEntryUpdated)

	ticketResponses, err := toTicketResponses(apiTickets)
	if err != nil {
		return nil, err
	}

	for _, ticketResponse := range ticketResponses {
		for playbookID := range ticketResponse.Playbooks {
			if err := runRootTask(ticketResponse, playbookID, db); err != nil {
				return nil, err
			}
		}
	}

	return ticketResponses, nil
}

func (db *Database) IndexRebuild(ctx context.Context) error {
	if err := db.Index.Truncate(); err != nil {
		return err
	}

	tickets, _, err := db.ticketListQuery(ctx, "FOR d IN @@collection RETURN d", nil, busdb.ReadOperation)
	if err != nil {
		return err
	}

	return batchIndex(db.Index, tickets)
}

func batchIndex(index *index.Index, tickets []*models.TicketSimpleResponse) error {
	var wg sync.WaitGroup
	var batch []*models.TicketSimpleResponse
	for _, ticket := range tickets {
		batch = append(batch, ticket)

		if len(batch) > 100 {
			wg.Add(1)
			go func(docs []*models.TicketSimpleResponse) {
				index.Index(docs)
				wg.Done()
			}(batch)
			batch = []*models.TicketSimpleResponse{}
		}
	}
	wg.Wait()
	return nil
}

func (db *Database) TicketGet(ctx context.Context, ticketID int64) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketReadFilter(ctx)
	if err != nil {
		return nil, err
	}

	return db.ticketGetQuery(ctx, ticketID, `LET d = DOCUMENT(@@collection, @ID) `+ticketFilterQuery+` RETURN d`, ticketFilterVars, busdb.ReadOperation)
}

func (db *Database) ticketGetQuery(ctx context.Context, ticketID int64, query string, bindVars map[string]interface{}, operation *busdb.Operation) (*models.TicketWithTickets, error) {
	if bindVars == nil {
		bindVars = map[string]interface{}{}
	}
	bindVars["@collection"] = TicketCollectionName
	if ticketID != 0 {
		bindVars["ID"] = fmt.Sprint(ticketID)
	}

	cur, _, err := db.Query(ctx, query, bindVars, operation)
	if err != nil {
		return nil, err
	}
	defer cur.Close()

	ticket := models.Ticket{}
	meta, err := cur.ReadDocument(ctx, &ticket)
	if err != nil {
		return nil, err
	}

	ticketSimpleResponse, err := toTicketSimpleResponse(meta.Key, &ticket)
	if err != nil {
		return nil, err
	}

	// index
	go db.Index.Index([]*models.TicketSimpleResponse{ticketSimpleResponse})

	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketReadFilter(ctx)
	if err != nil {
		return nil, err
	}

	// tickets
	ticketsQuery := `FOR vertex, edge IN OUTBOUND 
	DOCUMENT(@@tickets, @ID)
	GRAPH @graph
	FILTER IS_SAME_COLLECTION(@@collection, vertex)
	FILTER vertex != null
	LET d = DOCUMENT(@@collection, edge["_to"])
	` + ticketFilterQuery + `
	RETURN d`

	outTickets, _, err := db.ticketListQuery(ctx, ticketsQuery, mergeMaps(map[string]interface{}{
		"ID":       fmt.Sprint(ticketID),
		"graph":    TicketArtifactsGraphName,
		"@tickets": TicketCollectionName,
	}, ticketFilterVars), busdb.ReadOperation)
	if err != nil {
		return nil, err
	}

	ticketsQuery = `FOR vertex, edge IN INBOUND 
	DOCUMENT(@@tickets, @ID)
	GRAPH @graph
	FILTER IS_SAME_COLLECTION(@@collection, vertex)
	FILTER vertex != null
	LET d = DOCUMENT(@@collection, edge["_from"])
	` + ticketFilterQuery + `
	RETURN d`

	inTickets, _, err := db.ticketListQuery(ctx, ticketsQuery, mergeMaps(map[string]interface{}{
		"ID":       fmt.Sprint(ticketID),
		"graph":    TicketArtifactsGraphName,
		"@tickets": TicketCollectionName,
	}, ticketFilterVars), busdb.ReadOperation)
	if err != nil {
		return nil, err
	}

	var artifactNames []string
	for _, artifact := range ticketSimpleResponse.Artifacts {
		artifactNames = append(artifactNames, artifact.Name)
	}
	ticketsQuery = `FOR d IN @@collection
	FILTER d._key != @ID
	` + ticketFilterQuery + `
	FOR a IN NOT_NULL(d.artifacts, [])
	FILTER POSITION(@artifacts, a.name)
	RETURN d`
	sameArtifactTickets, _, err := db.ticketListQuery(ctx, ticketsQuery, mergeMaps(map[string]interface{}{
		"ID":        fmt.Sprint(ticketID),
		"artifacts": artifactNames,
	}, ticketFilterVars), busdb.ReadOperation)
	if err != nil {
		return nil, err
	}

	tickets := append(outTickets, inTickets...)
	tickets = append(tickets, sameArtifactTickets...)
	sort.Slice(tickets, func(i, j int) bool {
		return tickets[i].ID < tickets[j].ID
	})

	ticketResponse, err := toTicketResponse(ticketSimpleResponse)
	if err != nil {
		return nil, err
	}

	return toTicketWithTickets(ticketResponse, tickets), nil
}

func (db *Database) TicketUpdate(ctx context.Context, ticketID int64, ticket *models.Ticket) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	REPLACE d WITH @ticket IN @@collection
	RETURN NEW`
	ticket.Modified = time.Now().UTC() // TODO make setable?
	return db.ticketGetQuery(ctx, ticketID, query, mergeMaps(map[string]interface{}{"ticket": ticket}, ticketFilterVars), &busdb.Operation{
		Type: bus.DatabaseEntryUpdated, Ids: []driver.DocumentID{
			driver.NewDocumentID(TicketCollectionName, strconv.FormatInt(ticketID, 10)),
		},
	})
}

func (db *Database) TicketDelete(ctx context.Context, ticketID int64) error {
	_, err := db.TicketGet(ctx, ticketID)
	if err != nil {
		return err
	}

	_, err = db.ticketCollection.RemoveDocument(ctx, strconv.FormatInt(ticketID, 10))
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) TicketList(ctx context.Context, ticketType string, query string, sorts []string, desc []bool, offset, count int64) (*models.TicketList, error) {
	binVars := map[string]interface{}{}

	parser := &caql.Parser{Searcher: db.Index, Prefix: "d."}

	var typeString = ""
	if ticketType != "" {
		typeString = "FILTER d.type == @type "
		binVars["type"] = ticketType
	}

	var filterString = ""
	if query != "" {
		queryTree, err := parser.Parse(query)
		if err != nil {
			return nil, errors.New("invalid filter query: syntax error")
		}
		filterString, err = queryTree.String()
		if err != nil {
			return nil, fmt.Errorf("invalid filter query: %w", err)
		}
		filterString = "FILTER " + filterString
	}

	documentCount, err := db.TicketCount(ctx, typeString, filterString, binVars)
	if err != nil {
		return nil, err
	}

	sortQ := sortQuery(sorts, desc, binVars)
	binVars["offset"] = offset
	binVars["count"] = count

	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketReadFilter(ctx)
	if err != nil {
		return nil, err
	}

	q := `FOR d IN @@collection
		` + ticketFilterQuery + `
		` + sortQ + ` 
		` + typeString + `
		` + filterString + `
		LIMIT @offset, @count
		SORT d._key ASC
		RETURN d`
	// RETURN KEEP(d, "_key", "id", "name", "type", "created")`
	ticketList, _, err := db.ticketListQuery(ctx, q, mergeMaps(binVars, ticketFilterVars), busdb.ReadOperation)
	return &models.TicketList{
		Count:   documentCount,
		Tickets: ticketList,
	}, err
	// return map[string]interface{}{"tickets": ticketList, "count": documentCount}, err
}

func (db *Database) ticketListQuery(ctx context.Context, query string, bindVars map[string]interface{}, operation *busdb.Operation) ([]*models.TicketSimpleResponse, *models.LogEntry, error) {
	if bindVars == nil {
		bindVars = map[string]interface{}{}
	}
	bindVars["@collection"] = TicketCollectionName

	cursor, logEntry, err := db.Query(ctx, query, bindVars, operation)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close()

	var docs []*models.TicketSimpleResponse
	for {
		doc := models.Ticket{}
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, nil, err
		}

		resp, err := toTicketSimpleResponse(meta.Key, &doc)
		if err != nil {
			return nil, nil, err
		}

		docs = append(docs, resp)
	}

	return docs, logEntry, nil
}

func (db *Database) TicketCount(ctx context.Context, typequery, filterquery string, bindVars map[string]interface{}) (int, error) {
	if bindVars == nil {
		bindVars = map[string]interface{}{}
	}
	bindVars["@collection"] = TicketCollectionName

	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketReadFilter(ctx)
	if err != nil {
		return 0, err
	}

	countQuery := `RETURN LENGTH(FOR d IN @@collection ` + ticketFilterQuery + "  " + typequery + " " + filterquery + ` RETURN 1)`
	cursor, _, err := db.Query(ctx, countQuery, mergeMaps(bindVars, ticketFilterVars), busdb.ReadOperation)
	if err != nil {
		return 0, err
	}
	documentCount := 0
	_, err = cursor.ReadDocument(ctx, &documentCount)
	if err != nil {
		return 0, err
	}
	cursor.Close()
	return documentCount, nil
}

func sortQuery(paramsSort []string, paramsDesc []bool, bindVars map[string]interface{}) string {
	sort := ""
	if len(paramsSort) > 0 {
		var sorts []string
		for i, column := range paramsSort {
			colsort := fmt.Sprintf("d.@column%d", i)
			if len(paramsDesc) > i && paramsDesc[i] {
				colsort += " DESC"
			}
			sorts = append(sorts, colsort)
			bindVars[fmt.Sprintf("column%d", i)] = column
		}
		sort = "SORT " + strings.Join(sorts, ", ")
	}
	return sort
}

func mergeMaps(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	merged := map[string]interface{}{}
	for k, v := range a {
		merged[k] = v
	}
	for k, v := range b {
		merged[k] = v
	}
	return merged
}

func validate(e interface{}, schema *gojsonschema.Schema) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	res, err := schema.Validate(gojsonschema.NewStringLoader(string(b)))
	if err != nil {
		return err
	}

	if len(res.Errors()) > 0 {
		var l []string
		for _, e := range res.Errors() {
			l = append(l, e.String())
		}
		return fmt.Errorf("validation failed: %v", strings.Join(l, ", "))
	}
	return nil
}
