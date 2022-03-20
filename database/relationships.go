package database

import (
	"context"
	"errors"
	"strconv"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database/busdb"
)

func (db *Database) RelatedCreate(ctx context.Context, id, id2 int64) error {
	if id == id2 {
		return errors.New("tickets cannot relate to themself")
	}

	_, err := db.relatedCollection.CreateEdge(ctx, ctx, &driver.EdgeDocument{
		From: driver.DocumentID(TicketCollectionName + "/" + strconv.Itoa(int(id))),
		To:   driver.DocumentID(TicketCollectionName + "/" + strconv.Itoa(int(id2))),
	})

	return err
}

func (db *Database) RelatedBatchCreate(ctx context.Context, edges []*driver.EdgeDocument) error {
	_, err := db.relatedCollection.CreateEdges(ctx, edges)

	return err
}

func (db *Database) RelatedRemove(ctx context.Context, id, id2 int64) error {
	q := `
	FOR d in @@collection
	FILTER (d._from == @id && d._to == @id2) || (d._to == @id && d._from == @id2)
	REMOVE d in @@collection`
	_, _, err := db.Query(ctx, q, map[string]any{
		"@collection": RelatedTicketsCollectionName,
		"id":          driver.DocumentID(TicketCollectionName + "/" + strconv.Itoa(int(id))),
		"id2":         driver.DocumentID(TicketCollectionName + "/" + strconv.Itoa(int(id2))),
	}, &busdb.Operation{
		Type: bus.DatabaseEntryUpdated,
		Ids: []driver.DocumentID{
			driver.DocumentID(TicketCollectionName + "/" + strconv.Itoa(int(id))),
			driver.DocumentID(TicketCollectionName + "/" + strconv.Itoa(int(id2))),
		},
	})

	return err
}
