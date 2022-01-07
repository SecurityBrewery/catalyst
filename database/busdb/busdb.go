package busdb

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

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

type Operation struct {
	Type bus.DatabaseUpdateType
	Ids  []driver.DocumentID
}

var CreateOperation = &Operation{Type: bus.DatabaseEntryCreated}
var ReadOperation = &Operation{Type: bus.DatabaseEntryRead}

func (db BusDatabase) Query(ctx context.Context, query string, vars map[string]interface{}, operation *Operation) (cur driver.Cursor, logs *model.LogEntry, err error) {
	defer func() { err = toHTTPErr(err) }()

	cur, err = db.internal.Query(ctx, query, vars)
	if err != nil {
		return nil, nil, err
	}

	switch {
	case operation.Type == bus.DatabaseEntryCreated, operation.Type == bus.DatabaseEntryUpdated:
		if err := db.bus.PublishDatabaseUpdate(operation.Ids, operation.Type); err != nil {
			return nil, nil, err
		}
	}

	return cur, logs, err
}

func (db BusDatabase) Remove(ctx context.Context) (err error) {
	defer func() { err = toHTTPErr(err) }()

	return db.internal.Remove(ctx)
}

func (db BusDatabase) Collection(ctx context.Context, name string) (col driver.Collection, err error) {
	defer func() { err = toHTTPErr(err) }()

	return db.internal.Collection(ctx, name)
}

type Collection struct {
	internal driver.Collection
	db       *BusDatabase
}

func NewCollection(internal driver.Collection, db *BusDatabase) *Collection {
	return &Collection{internal: internal, db: db}
}

func (c Collection) CreateDocument(ctx, newctx context.Context, key string, document interface{}) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.CreateDocument(newctx, &Keyed{Key: key, Doc: document})
	if err != nil {
		return meta, err
	}

	err = c.db.bus.PublishDatabaseUpdate([]driver.DocumentID{meta.ID}, bus.DatabaseEntryCreated)
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func (c Collection) CreateEdge(ctx, newctx context.Context, edge *driver.EdgeDocument) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.CreateDocument(newctx, edge)
	if err != nil {
		return meta, err
	}

	err = c.db.bus.PublishDatabaseUpdate([]driver.DocumentID{meta.ID}, bus.DatabaseEntryCreated)
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func (c Collection) CreateEdges(ctx context.Context, edges []*driver.EdgeDocument) (meta driver.DocumentMetaSlice, err error) {
	defer func() { err = toHTTPErr(err) }()

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

	err = c.db.bus.PublishDatabaseUpdate(ids, bus.DatabaseEntryCreated)
	if err != nil {
		return metas, err
	}

	return metas, nil
}

func (c Collection) DocumentExists(ctx context.Context, id string) (exists bool, err error) {
	defer func() { err = toHTTPErr(err) }()

	return c.internal.DocumentExists(ctx, id)
}

func (c Collection) ReadDocument(ctx context.Context, key string, result interface{}) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.ReadDocument(ctx, key, result)

	return
}

func (c Collection) UpdateDocument(ctx context.Context, key string, update interface{}) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.UpdateDocument(ctx, key, update)
	if err != nil {
		return meta, err
	}

	return meta, c.db.bus.PublishDatabaseUpdate([]driver.DocumentID{meta.ID}, bus.DatabaseEntryUpdated)
}

func (c Collection) ReplaceDocument(ctx context.Context, key string, document interface{}) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.ReplaceDocument(ctx, key, document)
	if err != nil {
		return meta, err
	}

	return meta, c.db.bus.PublishDatabaseUpdate([]driver.DocumentID{meta.ID}, bus.DatabaseEntryUpdated)
}

func (c Collection) RemoveDocument(ctx context.Context, formatInt string) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	return c.internal.RemoveDocument(ctx, formatInt)
}

func (c Collection) Truncate(ctx context.Context) (err error) {
	defer func() { err = toHTTPErr(err) }()

	return c.internal.Truncate(ctx)
}

func toHTTPErr(err error) error {
	if err != nil {
		ae := driver.ArangoError{}
		if errors.As(err, &ae) {
			return &api.HTTPError{Status: ae.Code, Internal: err}
		}
	}
	return nil
}
