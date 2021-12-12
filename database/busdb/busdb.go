package busdb

import (
	"context"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/generated/models"
)

type Hook interface {
	PublishAction(action string, context, msg map[string]interface{}) error
	PublishUpdate(col, id string) error
}

// BusDatabase
//  1. Save entry to log
//  2. Send update ticket to bus
//  3. Add document to index
type BusDatabase struct {
	internal      driver.Database
	logCollection driver.Collection
	bus           *bus.Bus
	// index         *index.Index
}

func NewDatabase(ctx context.Context, internal driver.Database, b *bus.Bus) (*BusDatabase, error) {
	logCollection, err := internal.Collection(ctx, LogCollectionName)
	if err != nil {
		return nil, err
	}

	return &BusDatabase{
		internal:      internal,
		logCollection: logCollection,
		bus:           b,
	}, nil
}

type OperationType int

const (
	Create OperationType = iota
	Read                 = iota
	Update               = iota
)

type Operation struct {
	OperationType OperationType
	Ids           []driver.DocumentID
	Msg           string
}

var CreateOperation = &Operation{OperationType: Create}
var ReadOperation = &Operation{OperationType: Read}

func (db BusDatabase) Query(ctx context.Context, query string, vars map[string]interface{}, operation *Operation) (driver.Cursor, *models.LogEntry, error) {
	cur, err := db.internal.Query(ctx, query, vars)
	if err != nil {
		return nil, nil, err
	}

	var logs *models.LogEntry

	switch {
	case operation.OperationType == Update:
		if err := db.LogAndNotify(ctx, operation.Ids, operation.Msg); err != nil {
			return nil, nil, err
		}
	}

	return cur, logs, err
}

func (db BusDatabase) LogAndNotify(ctx context.Context, ids []driver.DocumentID, msg string) error {
	var logEntries []*models.LogEntry
	for _, i := range ids {
		logEntries = append(logEntries, &models.LogEntry{Reference: i.String(), Message: msg})
	}

	if err := db.LogBatchCreate(ctx, logEntries); err != nil {
		return err
	}

	return db.bus.PublishUpdate(ids)
}

func (db BusDatabase) Remove(ctx context.Context) error {
	return db.internal.Remove(ctx)
}

func (db BusDatabase) Collection(ctx context.Context, name string) (driver.Collection, error) {
	return db.internal.Collection(ctx, name)
}

type Collection struct {
	internal driver.Collection
	db       *BusDatabase
}

func NewCollection(internal driver.Collection, db *BusDatabase) *Collection {
	return &Collection{internal: internal, db: db}
}

func (c Collection) CreateDocument(ctx, newctx context.Context, key string, document interface{}) (driver.DocumentMeta, error) {
	meta, err := c.internal.CreateDocument(newctx, &Keyed{Key: key, Doc: document})
	if err != nil {
		return meta, err
	}

	err = c.db.LogAndNotify(ctx, []driver.DocumentID{meta.ID}, "Document created")
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func (c Collection) CreateEdge(ctx, newctx context.Context, edge *driver.EdgeDocument) (driver.DocumentMeta, error) {
	meta, err := c.internal.CreateDocument(newctx, edge)
	if err != nil {
		return meta, err
	}

	err = c.db.LogAndNotify(ctx, []driver.DocumentID{meta.ID}, "Document created")
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func (c Collection) CreateEdges(ctx context.Context, edges []*driver.EdgeDocument) (driver.DocumentMetaSlice, error) {
	metas, errs, err := c.internal.CreateDocuments(ctx, edges)
	if err != nil {
		return nil, err
	}
	if errs.FirstNonNil() != nil {
		return nil, errs.FirstNonNil()
	}

	var ids []driver.DocumentID
	for _, meta := range metas {
		ids = append(ids, meta.ID)
	}

	err = c.db.LogAndNotify(ctx, ids, "Document created")
	if err != nil {
		return metas, err
	}

	return metas, nil
}

func (c Collection) DocumentExists(ctx context.Context, id string) (bool, error) {
	return c.internal.DocumentExists(ctx, id)
}

func (c Collection) ReadDocument(ctx context.Context, key string, result interface{}) (driver.DocumentMeta, error) {
	return c.internal.ReadDocument(ctx, key, result)
}

func (c Collection) UpdateDocument(ctx context.Context, key string, update interface{}) (driver.DocumentMeta, error) {
	meta, err := c.internal.UpdateDocument(ctx, key, update)
	if err != nil {
		return meta, err
	}

	return meta, c.db.bus.PublishUpdate([]driver.DocumentID{meta.ID})
}

func (c Collection) ReplaceDocument(ctx context.Context, key string, document interface{}) (driver.DocumentMeta, error) {
	meta, err := c.internal.ReplaceDocument(ctx, key, document)
	if err != nil {
		return meta, err
	}

	return meta, c.db.bus.PublishUpdate([]driver.DocumentID{meta.ID})
}

func (c Collection) RemoveDocument(ctx context.Context, formatInt string) (driver.DocumentMeta, error) {
	return c.internal.RemoveDocument(ctx, formatInt)
}
