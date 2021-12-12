package database

import (
	"context"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
)

func (db *Database) ArtifactGet(ctx context.Context, id int64, name string) (*models.Artifact, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID) 
	` + ticketFilterQuery + `
	FOR a in NOT_NULL(d.artifacts, [])
	FILTER a.name == @name
	RETURN a`
	cursor, _, err := db.Query(ctx, query, mergeMaps(ticketFilterVars, map[string]interface{}{
		"@collection": TicketCollectionName,
		"ID":          fmt.Sprint(id),
		"name":        name,
	}), busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var doc models.Artifact
	_, err = cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (db *Database) ArtifactUpdate(ctx context.Context, id int64, name string, artifact *models.Artifact) (*models.TicketWithTickets, error) {
	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	FOR a IN NOT_NULL(d.artifacts, [])
	FILTER a.name == @name
	LET newartifacts = APPEND(REMOVE_VALUE(d.artifacts, a), @artifact)
	UPDATE d WITH { "artifacts": newartifacts } IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{
		"@collection": TicketCollectionName,
		"ID":          id,
		"name":        name,
		"artifact":    artifact,
	}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.DocumentID(fmt.Sprintf("%s/%d", TicketCollectionName, id)),
		},
		Msg: fmt.Sprintf("Update artifact %s", name),
	})
}

func (db *Database) EnrichArtifact(ctx context.Context, id int64, name string, enrichmentForm *models.EnrichmentForm) (*models.TicketWithTickets, error) {
	enrichment := models.Enrichment{time.Now().UTC(), enrichmentForm.Data, enrichmentForm.Name}

	ticketFilterQuery, ticketFilterVars, err := db.Hooks.TicketWriteFilter(ctx)
	if err != nil {
		return nil, err
	}

	query := `LET d = DOCUMENT(@@collection, @ID)
	` + ticketFilterQuery + `
	FOR a IN d.artifacts
	FILTER a.name == @name
	LET enrichments = NOT_NULL(a.enrichments, {})
	LET newenrichments = MERGE(enrichments, ZIP( [@enrichmentname], [@enrichment]) )
	LET newartifacts = APPEND(REMOVE_VALUE(d.artifacts, a), MERGE(a, { "enrichments": newenrichments }))
	UPDATE d WITH { "artifacts": newartifacts } IN @@collection
	RETURN NEW`
	return db.ticketGetQuery(ctx, id, query, mergeMaps(map[string]interface{}{
		"@collection":    TicketCollectionName,
		"ID":             id,
		"name":           name,
		"enrichmentname": enrichment.Name,
		"enrichment":     enrichment,
	}, ticketFilterVars), &busdb.Operation{
		OperationType: busdb.Update,
		Ids: []driver.DocumentID{
			driver.DocumentID(fmt.Sprintf("%s/%d", TicketCollectionName, id)),
		},
		Msg: fmt.Sprintf("Run %s on artifact", enrichment.Name),
	})
}
